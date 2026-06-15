package seeders

import "database/sql"

func SeedRoles(db *sql.DB) error {
	roles := []struct{ name, displayName, description string }{
		{"super_admin", "Super Administrator", "Full unrestricted access"},
		{"admin",       "Administrator",       "Manages users and content"},
		{"user",        "User",                "Standard application user"},
	}

	for _, r := range roles {
		_, err := db.Exec(
			`INSERT IGNORE INTO roles (name, display_name, description) VALUES (?, ?, ?)`,
			r.name, r.displayName, r.description,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
