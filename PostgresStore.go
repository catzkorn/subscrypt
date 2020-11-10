package main

import (
	"database/sql"
	"fmt"
)

type PostgresStore struct {
	database *sql.DB
}

// NewPostgresStore starts connection with database
func NewPostgresStore(databaseDSN string) (*PostgresStore, error) {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unexpected connection error: %w", err)
	}

	return &PostgresStore{database: db}, nil
}
