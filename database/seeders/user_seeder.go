package seeders

import (
	"context"
	"database/sql"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func SeedUsers(db *sql.DB) error {
	ctx := context.Background()

	var count int
	_ = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if count > 0 {
		slog.Info("users table already seeded — skipping")
		return nil
	}

	seeds := []struct {
		name, email, password, role string
	}{
		{"Super Admin", "superadmin@vuelang.dev", "password123", "super_admin"},
		{"Admin User",  "admin@vuelang.dev",      "password123", "admin"},
		{"Demo User",   "demo@vuelang.dev",        "password123", "user"},
	}

	for _, s := range seeds {
		hashed, err := bcrypt.GenerateFromPassword([]byte(s.password), 12)
		if err != nil {
			return err
		}

		res, err := db.ExecContext(ctx,
			`INSERT INTO users (name, email, password, is_active) VALUES (?, ?, ?, 1)`,
			s.name, s.email, string(hashed),
		)
		if err != nil {
			return err
		}
		userID, _ := res.LastInsertId()

		// Assign role
		var roleID uint
		_ = db.QueryRowContext(ctx, `SELECT id FROM roles WHERE name = ?`, s.role).Scan(&roleID)
		if roleID > 0 {
			_, _ = db.ExecContext(ctx,
				`INSERT IGNORE INTO role_user (user_id, role_id) VALUES (?, ?)`,
				userID, roleID,
			)
		}
		slog.Info("seeded user", "email", s.email, "role", s.role)
	}
	return nil
}
