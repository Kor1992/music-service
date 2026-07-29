package repository

import (
	"context"
	"errors"
	"music/internal/domain"
)

var ErrNotFound = errors.New("not found")

type TrackRepository interface {
	Create(ctx context.Context, title, prompt string, userID int) (*domain.Track, error)
	GetByID(ctx context.Context, id int) (*domain.Track, error)
	List(ctx context.Context) ([]*domain.Track, error)
	UpdateStatus(ctx context.Context, id int, status string, audioURL string) error
}
