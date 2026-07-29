run:
	go run cmd/server/main.go

db-up:
	docker compose up -d db

db-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/musicdb?sslmode=disable" up
migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/musicdb?sslmode=disable" down
