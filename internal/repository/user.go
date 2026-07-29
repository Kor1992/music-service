package repository

import (
	"context"
	"music/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, email, passwordHash, role string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	UpdateSubscription(ctx context.Context, id int, subscription string) error
}
