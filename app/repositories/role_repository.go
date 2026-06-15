package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"vuelang/app/models"
)

// RoleRepository handles role and role-user assignment queries.
type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) All(ctx context.Context) ([]models.Role, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, display_name, description, created_at FROM roles ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("RoleRepository.All: %w", err)
	}
	defer rows.Close()

	var list []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.DisplayName, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, role)
	}
	return list, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, display_name, description, created_at FROM roles WHERE name = ?`, name,
	).Scan(&role.ID, &role.Name, &role.DisplayName, &role.Description, &role.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &role, err
}

// GetUserRoles returns the role names assigned to a user.
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT ro.name FROM roles ro
		 INNER JOIN role_user ru ON ru.role_id = ro.id
		 WHERE ru.user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("RoleRepository.GetUserRoles: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

// AssignRole assigns a named role to a user.
func (r *RoleRepository) AssignRole(ctx context.Context, userID uint, roleName string) error {
	role, err := r.FindByName(ctx, roleName)
	if err != nil {
		return err
	}
	if role == nil {
		return fmt.Errorf("role %q not found", roleName)
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT IGNORE INTO role_user (user_id, role_id) VALUES (?, ?)`,
		userID, role.ID,
	)
	return err
}

// RemoveRole removes a named role from a user.
func (r *RoleRepository) RemoveRole(ctx context.Context, userID uint, roleName string) error {
	role, err := r.FindByName(ctx, roleName)
	if err != nil || role == nil {
		return err
	}
	_, err = r.db.ExecContext(ctx,
		`DELETE FROM role_user WHERE user_id = ? AND role_id = ?`, userID, role.ID)
	return err
}
