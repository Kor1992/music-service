package postgres

import (
	"context"
	"errors"
	"music/internal/domain"
	"music/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QueueRepo struct {
	pool *pgxpool.Pool
}

func NewQueueRepo(pool *pgxpool.Pool) repository.QueueRepository {
	return &QueueRepo{pool: pool}

}

func (r *QueueRepo) Enqueue(ctx context.Context, trackID int) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO music.generation_queue (track_id) VALUES ($1)",
		trackID,
	)
	return err
}

func (r *QueueRepo) Dequeue(ctx context.Context) (*domain.QueueItem, error) {
	item := &domain.QueueItem{}
	err := r.pool.QueryRow(ctx,
		`UPDATE music.generation_queue 
         SET status = 'processing', attempts = attempts + 1, updated_at = NOW()
         WHERE id = (
             SELECT id FROM music.generation_queue 
             WHERE status = 'pending' 
             ORDER BY created_at 
             LIMIT 1 
             FOR UPDATE SKIP LOCKED
         )
         RETURNING id, track_id, status, attempts, created_at, updated_at`,
	).Scan(&item.ID, &item.TrackID, &item.Status, &item.Attempts, &item.CreatedAt, &item.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *QueueRepo) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE music.generation_queue SET status = $1, updated_at = NOW() WHERE id = $2",
		status, id,
	)
	return err
}
