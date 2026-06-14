package migrations

import "database/sql"

func createStoresTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS stores (
			id         BIGINT UNSIGNED                    AUTO_INCREMENT PRIMARY KEY,
			name       VARCHAR(150)                       NOT NULL,
			code       VARCHAR(20)                        NOT NULL UNIQUE,
			type       ENUM('branch','warehouse')         NOT NULL DEFAULT 'branch',
			phone      VARCHAR(30)                        NOT NULL DEFAULT '',
			email      VARCHAR(100)                       NOT NULL DEFAULT '',
			address    TEXT                               NOT NULL DEFAULT '',
			city       VARCHAR(100)                       NOT NULL DEFAULT '',
			is_active  TINYINT(1)                        NOT NULL DEFAULT 1,
			created_at DATETIME                           NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME                           NOT NULL DEFAULT CURRENT_TIMESTAMP
			                                                       ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_code      (code),
			INDEX idx_type      (type),
			INDEX idx_is_active (is_active)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
