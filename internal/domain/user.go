package domain

import "time"

type User struct {
	ID           int        `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	Subscription string     `json:"subscription"`
	TrialEndsAt  *time.Time `json:"trial_ends_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}
