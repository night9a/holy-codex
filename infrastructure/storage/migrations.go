package storage

import "database/sql"

// RunMigrations creates schema if it doesn't exist, applying each version in order.
func RunMigrations(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	for _, m := range migrations {
		applied, err := isMigrationApplied(db, m.version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if _, err := db.Exec(m.sql); err != nil {
			return err
		}
		if err := recordMigration(db, m.version); err != nil {
			return err
		}
	}
	return nil
}

// migration holds a single versioned DDL statement.
type migration struct {
	version int
	sql     string
}

var migrations = []migration{
	{
		version: 1,
		sql: `
			CREATE TABLE IF NOT EXISTS entries (
				id         TEXT PRIMARY KEY,
				user_id    TEXT NOT NULL DEFAULT 'default',
				title      TEXT NOT NULL DEFAULT '',
				body       TEXT NOT NULL DEFAULT '',
				tags       TEXT NOT NULL DEFAULT '[]',
				mood       TEXT NOT NULL DEFAULT 'calm',
				created_at DATETIME NOT NULL,
				updated_at DATETIME NOT NULL,
				is_synced  INTEGER NOT NULL DEFAULT 0
			);
			CREATE INDEX IF NOT EXISTS idx_entries_user ON entries(user_id);
			CREATE INDEX IF NOT EXISTS idx_entries_created ON entries(created_at DESC);
		`,
	},
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func isMigrationApplied(db *sql.DB, version int) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version=?`, version).Scan(&count)
	return count > 0, err
}

func recordMigration(db *sql.DB, version int) error {
	_, err := db.Exec(`INSERT INTO schema_migrations(version) VALUES(?)`, version)
	return err
}