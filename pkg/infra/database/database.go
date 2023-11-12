package database

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func NewDatabase(databaseURL string) (*sql.DB, *bun.DB, error) {
	sqlDB, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()

		return nil, nil, err
	}

	return sqlDB, bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns()), nil
}
