package migrations

import "database/sql"

func createRolesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			name         VARCHAR(50)     NOT NULL UNIQUE,
			display_name VARCHAR(100)    NOT NULL DEFAULT '',
			description  VARCHAR(255)    NOT NULL DEFAULT '',
			created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
