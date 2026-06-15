package models

import "time"

// PasswordReset maps to the `password_resets` table.
type PasswordReset struct {
	Email     string    `json:"email"`
	Token     string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
