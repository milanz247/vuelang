package models

import "time"

// User maps to the `users` table.
// The Password field is always omitted from JSON serialisation.
type User struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	Email             string     `json:"email"`
	Password          string     `json:"-"`
	EmailVerifiedAt   *time.Time `json:"email_verified_at"`
	IsActive          bool       `json:"is_active"`
	Roles             []string   `json:"roles,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
