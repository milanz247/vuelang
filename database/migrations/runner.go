// Package migrations runs numbered schema changes in order.
// Every migration uses CREATE TABLE IF NOT EXISTS / ALTER TABLE IF NOT EXISTS
// so it is safe to re-run on every boot.
//
// To add a migration:
//  1. Create  database/migrations/0007_create_products_table.go
//  2. Add an entry to the list below.
package migrations

import (
	"database/sql"
	"fmt"
	"log/slog"
)

type migration struct {
	name string
	up   func(*sql.DB) error
}

func Run(db *sql.DB) error {
	list := []migration{
		{"0001_create_users_table",          createUsersTable},
		{"0002_create_roles_table",          createRolesTable},
		{"0003_create_role_user_table",      createRoleUserTable},
		{"0004_create_password_resets_table", createPasswordResetsTable},
		{"0005_create_refresh_tokens_table", createRefreshTokensTable},
		{"0006_create_audit_logs_table",     createAuditLogsTable},
	}

	for _, m := range list {
		if err := m.up(db); err != nil {
			return fmt.Errorf("migration %s: %w", m.name, err)
		}
		slog.Debug("migrated: " + m.name)
	}
	slog.Info("all migrations applied")
	return nil
}
