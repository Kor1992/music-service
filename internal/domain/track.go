package domain

import "time"

type Track struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Prompt    string    `json:"prompt"`
	AudioURL  *string   `json:"audio_url,omitempty"`
	Status    string    `json:"status"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
