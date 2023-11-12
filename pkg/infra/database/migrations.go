package database

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func RunMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "./scripts/db/migrations",
	}

	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)

	return err
}
