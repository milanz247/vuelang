// Package migrations — one file per table, run in numbered order.
//
// To add a new table:
//   1. Create  database/migrations/0003_create_products_table.go
//   2. Write   createProductsTable(db *sql.DB) error  inside it
//   3. Add     {"0003_create_products_table", createProductsTable}  to the list below
//
// All migrations use CREATE TABLE IF NOT EXISTS — safe to re-run on every boot.
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

// Run executes every migration in the order listed.
func Run(db *sql.DB) error {
	list := []migration{
		{"0001_create_users_table",  createUsersTable},
		{"0002_create_stores_table", createStoresTable},
		// {"0003_create_products_table", createProductsTable},
	}

	for _, m := range list {
		if err := m.up(db); err != nil {
			return fmt.Errorf("migration %s: %w", m.name, err)
		}
		slog.Info("migrated: " + m.name)
	}
	return nil
}
