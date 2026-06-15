package models

import "time"

// RefreshToken maps to the `refresh_tokens` table.
type RefreshToken struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
