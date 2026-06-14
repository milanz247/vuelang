package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User maps to the `users` table.
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // never sent to client
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ── DB methods (like Laravel Eloquent static methods) ─────────────────────────

// All returns every user ordered newest first.
func UserAll(ctx context.Context, db *sql.DB) ([]User, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, name, email, is_active, created_at, updated_at
		 FROM users ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("User.All: %w", err)
	}
	defer rows.Close()

	var list []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email,
			&u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, nil
}

// Find returns one user by ID, or nil if not found.
func UserFind(ctx context.Context, db *sql.DB, id uint) (*User, error) {
	var u User
	err := db.QueryRowContext(ctx,
		`SELECT id, name, email, is_active, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// FindByEmail returns one user by email, or nil if not found.
func UserFindByEmail(ctx context.Context, db *sql.DB, email string) (*User, error) {
	var u User
	err := db.QueryRowContext(ctx,
		`SELECT id, name, email, password, is_active, created_at, updated_at
		 FROM users WHERE email = ?`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &u, err
}

// Create inserts a new user and returns it with the generated ID.
func UserCreate(ctx context.Context, db *sql.DB, name, email, password string) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt: %w", err)
	}
	res, err := db.ExecContext(ctx,
		`INSERT INTO users (name, email, password, is_active) VALUES (?, ?, ?, 1)`,
		name, email, string(hashed),
	)
	if err != nil {
		return nil, fmt.Errorf("User.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return UserFind(ctx, db, uint(id))
}

// Update saves name, email, is_active changes for an existing user.
func UserUpdate(ctx context.Context, db *sql.DB, id uint, name, email string, isActive bool) (*User, error) {
	_, err := db.ExecContext(ctx,
		`UPDATE users SET name = ?, email = ?, is_active = ? WHERE id = ?`,
		name, email, isActive, id,
	)
	if err != nil {
		return nil, fmt.Errorf("User.Update: %w", err)
	}
	return UserFind(ctx, db, id)
}

// Delete removes a user by ID.
func UserDelete(ctx context.Context, db *sql.DB, id uint) error {
	_, err := db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}
