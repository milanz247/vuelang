package models

import "time"

// AuditLog maps to the `audit_logs` table.
// Every state-changing operation should produce an entry.
type AuditLog struct {
	ID         uint      `json:"id"`
	UserID     *uint     `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID *uint     `json:"resource_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	StatusCode int       `json:"status_code"`
	CreatedAt  time.Time `json:"created_at"`
}
