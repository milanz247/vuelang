package migrations

import "database/sql"

func createPasswordResetsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS password_resets (
			email      VARCHAR(150) NOT NULL,
			token      VARCHAR(255) NOT NULL,
			expires_at DATETIME     NOT NULL,
			created_at DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (email),
			INDEX idx_token (token)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
