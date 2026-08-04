package repository

import (
	"context"
	"music/internal/domain"
)

type QueueRepository interface {
	Enqueue(ctx context.Context, trackID int) error
	Dequeue(ctx context.Context) (*domain.QueueItem, error)
	UpdateStatus(ctx context.Context, id int, status string) error
}
