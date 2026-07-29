package postgres

import (
	"context"
	"errors"
	"music/internal/domain"
	"music/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrackRepo struct {
	pool *pgxpool.Pool
}

func NewTrackRepo(pool *pgxpool.Pool) repository.TrackRepository {
	return &TrackRepo{pool: pool}
}

func (r *TrackRepo) Create(ctx context.Context, title, prompt string, userID int) (*domain.Track, error) {
	t := &domain.Track{}

	query := `
	INSERT INTO music.tracks (title, prompt, user_id)
	VALUES ($1, $2, $3)
	RETURNING id, title, prompt, audio_url, status, user_id, created_at
	`

	err := r.pool.QueryRow(ctx, query, title, prompt, userID).Scan(
		&t.ID,
		&t.Title,
		&t.Prompt,
		&t.AudioURL,
		&t.Status,
		&t.UserID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TrackRepo) GetByID(ctx context.Context, id int) (*domain.Track, error) {
	t := &domain.Track{}

	query := `
	SELECT id, title, prompt, audio_url, status, user_id, created_at
	FROM music.tracks
	WHERE id = $1
	`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&t.ID,
		&t.Title,
		&t.Prompt,
		&t.AudioURL,
		&t.Status,
		&t.UserID,
		&t.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}

	return t, nil
}

func (r *TrackRepo) List(ctx context.Context) ([]*domain.Track, error) {
	query := `
	SELECT id, title, prompt, audio_url, status, user_id, created_at
	FROM music.tracks
	ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tracks []*domain.Track

	for rows.Next() {
		t := &domain.Track{}
		if err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Prompt,
			&t.AudioURL,
			&t.Status,
			&t.UserID,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		tracks = append(tracks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tracks, nil
}

func (r *TrackRepo) UpdateStatus(ctx context.Context, id int, status string, audioURL string) error {
	query := `
	UPDATE music.tracks
	SET status = $1, audio_url = $2
	WHERE id = $3
	`

	_, err := r.pool.Exec(ctx, query, status, audioURL, id)

	return err
}
