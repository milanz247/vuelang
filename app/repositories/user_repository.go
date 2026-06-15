package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"vuelang/app/models"
)

// UserRepository handles all database operations for the users table.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

const userSelectCols = `u.id, u.name, u.email, u.password, u.email_verified_at, u.is_active, u.created_at, u.updated_at`

func scanUser(row interface{ Scan(...any) error }) (*models.User, error) {
	var u models.User
	if err := row.Scan(
		&u.ID, &u.Name, &u.Email, &u.Password,
		&u.EmailVerifiedAt, &u.IsActive, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) All(ctx context.Context) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT `+userSelectCols+` FROM users u ORDER BY u.id DESC`)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.All: %w", err)
	}
	defer rows.Close()

	var list []models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *u)
	}
	return list, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx,
		`SELECT `+userSelectCols+` FROM users u WHERE u.id = ?`, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx,
		`SELECT `+userSelectCols+` FROM users u WHERE u.email = ?`, email))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepository) Create(ctx context.Context, name, email, hashedPassword string) (*models.User, error) {
	res, err := r.db.ExecContext(ctx,
		`INSERT INTO users (name, email, password, is_active) VALUES (?, ?, ?, 1)`,
		name, email, hashedPassword,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.Create: %w", err)
	}
	id, _ := res.LastInsertId()
	return r.FindByID(ctx, uint(id))
}

func (r *UserRepository) Update(ctx context.Context, id uint, name, email string, isActive bool) (*models.User, error) {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET name = ?, email = ?, is_active = ? WHERE id = ?`,
		name, email, isActive, id,
	)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.Update: %w", err)
	}
	return r.FindByID(ctx, id)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password = ? WHERE id = ?`, hashedPassword, id)
	return err
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET email_verified_at = NOW() WHERE id = ?`, id)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	return err
}

// ── Refresh tokens ────────────────────────────────────────────────────────────

func (r *UserRepository) CreateRefreshToken(ctx context.Context, userID uint, token string, expiresAt interface{}) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)`,
		userID, token, expiresAt,
	)
	return err
}

func (r *UserRepository) FindRefreshToken(ctx context.Context, token string) (uint, bool, error) {
	var userID uint
	var expired bool
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id, expires_at < NOW() FROM refresh_tokens WHERE token = ?`, token,
	).Scan(&userID, &expired)
	if err == sql.ErrNoRows {
		return 0, false, nil
	}
	return userID, !expired, err
}

func (r *UserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE token = ?`, token)
	return err
}

func (r *UserRepository) DeleteUserRefreshTokens(ctx context.Context, userID uint) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE user_id = ?`, userID)
	return err
}

// ── Password resets ───────────────────────────────────────────────────────────

func (r *UserRepository) CreatePasswordReset(ctx context.Context, email, token string, expiresAt interface{}) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO password_resets (email, token, expires_at) VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE token = VALUES(token), expires_at = VALUES(expires_at), created_at = NOW()`,
		email, token, expiresAt,
	)
	return err
}

func (r *UserRepository) FindPasswordReset(ctx context.Context, token string) (string, bool, error) {
	var email string
	var expired bool
	err := r.db.QueryRowContext(ctx,
		`SELECT email, expires_at < NOW() FROM password_resets WHERE token = ?`, token,
	).Scan(&email, &expired)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	return email, !expired, err
}

func (r *UserRepository) DeletePasswordReset(ctx context.Context, email string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM password_resets WHERE email = ?`, email)
	return err
}

// ── Audit logs ────────────────────────────────────────────────────────────────

func (r *UserRepository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO audit_logs (user_id, action, resource, resource_id, ip_address, user_agent, status_code)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		log.UserID, log.Action, log.Resource, log.ResourceID,
		log.IPAddress, log.UserAgent, log.StatusCode,
	)
	return err
}
