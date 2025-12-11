package assignment

import (
	"database/sql"
)

type PostgresAssignmentStorage struct {
	db *sql.DB
}

// New creates a new PostgresAssignmentStorage instance.
// That used to interact with the assignments table.
func New(db *sql.DB) *PostgresAssignmentStorage {
	return &PostgresAssignmentStorage{db: db}
}
