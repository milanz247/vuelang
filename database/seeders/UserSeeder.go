// Package seeders populates the DB with starter data for development.
// Never run automatically in production.
//
// Usage (in main.go, only when ENV=development and --seed flag passed):
//
//	seeders.UserSeeder(db)
package seeders

import (
	"context"
	"database/sql"
	"log/slog"

	"go-cloud-erp/app/models"
)

// UserSeeder inserts demo users if the table is empty.
func UserSeeder(db *sql.DB) error {
	ctx := context.Background()

	users, err := models.UserAll(ctx, db)
	if err != nil {
		return err
	}
	if len(users) > 0 {
		slog.Info("users table already has data — skipping seeder")
		return nil
	}

	seeds := []struct{ name, email, password string }{
		{"Admin",     "admin@vuelang.dev", "password123"},
		{"Demo User", "demo@vuelang.dev",  "password123"},
	}

	for _, s := range seeds {
		u, err := models.UserCreate(ctx, db, s.name, s.email, s.password)
		if err != nil {
			return err
		}
		slog.Info("seeded user", "email", u.Email)
	}
	return nil
}
