package submission

import (
	"database/sql"
)

type PostgresSubmissionStorage struct {
	db *sql.DB
}

// New creates a new PostgresSubmissionStorage instance.
// That used to interact with the submissions table.
func New(db *sql.DB) *PostgresSubmissionStorage {
	return &PostgresSubmissionStorage{db: db}
}
