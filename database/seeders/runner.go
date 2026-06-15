// Package seeders populates the DB with starter data for development.
// Never run automatically in production.
package seeders

import (
	"database/sql"
	"log/slog"
)

func RunAll(db *sql.DB) error {
	runners := []struct {
		name string
		fn   func(*sql.DB) error
	}{
		{"roles",  SeedRoles},
		{"users",  SeedUsers},
	}

	for _, r := range runners {
		if err := r.fn(db); err != nil {
			return err
		}
		slog.Info("seeded: " + r.name)
	}
	return nil
}
