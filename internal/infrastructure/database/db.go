package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func New(dbstring string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dbstring)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
