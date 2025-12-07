package postgres

import (
	"database/sql"
	"fmt"

	"sso/internal/storage"

	_ "github.com/jackc/pgx/v5"
)

type Storage struct {
	db *sql.DB
	storage.UserStorage
	storage.RoleStorage
	storage.AppStorage
}

// New creates a new instance of PostgreSQL storage
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

	return &Storage{
		db:          db,
		UserStorage: NewUserStorage(db),
		RoleStorage: NewRoleStorage(db),
		AppStorage:  NewAppStorage(db),
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
