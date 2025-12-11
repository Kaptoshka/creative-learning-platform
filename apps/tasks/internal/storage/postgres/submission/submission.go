package submission

import (
	"context"
	"database/sql"

	"tasks/internal/domain/models"
)

type PostgresSubmissionStorage struct {
	db *sql.DB
}

// New creates a new PostgresSubmissionStorage instance.
// That used to interact with the submissions table.
func New(db *sql.DB) *PostgresSubmissionStorage {
	return &PostgresSubmissionStorage{db: db}
}

// TODO redo cause now it doesn't work correctly
func (s *PostgresSubmissionStorage) SubmissionByAssignmentID(
	ctx context.Context,
	assignmentID int64,
) (models.Submission, error) {
	const op = "storage.postgres.SubmissionByAssignmentID"

	query := `
		SELECT *
		FROM submissions
		WHERE assignment_id = $1
	`

}
