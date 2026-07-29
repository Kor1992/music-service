package postgres

import (
	"context"
	"errors"
	"music/internal/domain"
	"music/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) repository.UserRepository {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, role string) (*domain.User, error) {
	u := &domain.User{}
	query := `
INSERT INTO music.users (email, password_hash, role) 
VALUES ($1, $2, $3)
RETURNING id, email, role, subscription, trial_ends_at, created_at
`
	err := r.pool.QueryRow(ctx, query, email, passwordHash, role).Scan(
		&u.ID,
		&u.Email,
		&u.Role,
		&u.Subscription,
		&u.TrialEndsAt,
		&u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	query := `
	SELECT id, email, password_hash, role, subscription, trial_ends_at, created_at
	FROM music.users
	WHERE email = $1
	`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Subscription,
		&u.TrialEndsAt,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	u := &domain.User{}
	query := `
	SELECT id, email, password_hash, role, subscription, trial_ends_at, created_at
	FROM music.users
	WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.Role,
		&u.Subscription,
		&u.TrialEndsAt,
		&u.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return u, nil
}

func (r *UserRepo) UpdateSubscription(ctx context.Context, id int, subscription string) error {
	_, err := r.pool.Exec(ctx, "UPDATE music.users SET subscription = $1 WHERE id = $2", subscription, id)
	return err
}
