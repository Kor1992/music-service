package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"music/internal/handler"
	"music/internal/middleware"
	"music/internal/repository/postgres"
	"music/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to create pool:", err)
	}
	defer pool.Close()

	for i := 0; i < 10; i++ {
		if err := pool.Ping(ctx); err == nil {
			break
		}
		log.Printf("Waiting for database... attempt %d", i+1)
		time.Sleep(2 * time.Second)
	}
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Database is not ready:", err)
	}
	log.Println("Connected to database")

	// ---------- Инициализация ----------
	userRepo := postgres.NewUserRepo(pool)
	trackRepo := postgres.NewTrackRepo(pool)

	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	trackSvc := service.NewTrackService(trackRepo)
	trackHandler := handler.NewTrackHandler(trackSvc)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in .env")
	}

	// ---------- Роутер ----------
	mux := http.NewServeMux()

	// Публичные маршруты аутентификации
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)

	// Публичные маршруты треков
	mux.HandleFunc("GET /tracks", trackHandler.List)
	mux.HandleFunc("GET /tracks/{id}", trackHandler.GetById)

	// Защищённый маршрут создания трека (JWT + проверка подписки)
	mux.Handle("POST /tracks",
		middleware.AuthMiddleware(jwtSecret)(
			middleware.RequireSubscription(userRepo)(
				http.HandlerFunc(trackHandler.Create),
			),
		),
	)

	mux.Handle("POST /tracks/generate",
		middleware.AuthMiddleware(jwtSecret)(
			middleware.RequireSubscription(userRepo)(
				http.HandlerFunc(trackHandler.Generate),
			),
		),
	)

	mux.HandleFunc("GET /tracks/{id}/stream", trackHandler.Stream)

	mux.Handle("GET /tracks/{id}/stream",
		middleware.AuthMiddleware(jwtSecret)(
			http.HandlerFunc(trackHandler.Stream),
		),
	)

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fs)

	loggedMux := middleware.Logging(mux)

	// ---------- Старт ----------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatal(err)
	}
}
