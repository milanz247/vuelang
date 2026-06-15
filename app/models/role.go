package models

import "time"

// Role maps to the `roles` table.
type Role struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Built-in role names — use these constants everywhere instead of raw strings.
const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleUser       = "user"
)
