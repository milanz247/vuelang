package migrations

import "database/sql"

func createAuditLogsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS audit_logs (
			id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			user_id     BIGINT UNSIGNED NULL,
			action      VARCHAR(100)    NOT NULL,
			resource    VARCHAR(100)    NOT NULL DEFAULT '',
			resource_id BIGINT UNSIGNED NULL,
			ip_address  VARCHAR(45)     NOT NULL DEFAULT '',
			user_agent  VARCHAR(512)    NOT NULL DEFAULT '',
			status_code SMALLINT        NOT NULL DEFAULT 0,
			created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_user_id    (user_id),
			INDEX idx_action     (action),
			INDEX idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
	return err
}
