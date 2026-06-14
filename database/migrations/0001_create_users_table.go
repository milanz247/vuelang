package migrations

import "database/sql"

func createUsersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id         BIGINT UNSIGNED  AUTO_INCREMENT PRIMARY KEY,
			name       VARCHAR(100)     NOT NULL,
			email      VARCHAR(150)     NOT NULL UNIQUE,
			password   VARCHAR(255)     NOT NULL,
			is_active  TINYINT(1)      NOT NULL DEFAULT 1,
			created_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP
			                                    ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_email     (email),
			INDEX idx_is_active (is_active)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
