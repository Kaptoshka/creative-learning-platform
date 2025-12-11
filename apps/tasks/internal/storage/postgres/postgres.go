package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5"
)

type Storage struct {
	db *sql.DB
}

func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
