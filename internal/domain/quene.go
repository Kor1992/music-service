package domain

import "time"

type QueueItem struct {
	ID        int       `json:"id"`
	TrackID   int       `json:"track_id"`
	Status    string    `json:"status"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
