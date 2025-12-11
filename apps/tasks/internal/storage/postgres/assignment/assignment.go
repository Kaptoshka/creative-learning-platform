package assignment

import (
	"context"
	"database/sql"

	"tasks/internal/domain/models"
)

type PostgresAssignmentStorage struct {
	db *sql.DB
}

// New creates a new PostgresAssignmentStorage instance.
// That used to interact with the assignments table.
func New(db *sql.DB) *PostgresAssignmentStorage {
	return &PostgresAssignmentStorage{db: db}
}

func (s *PostgresAssignmentStorage) Assignment(
	ctx context.Context,
	assignmentID int64,
) (models.Assignment, error) {
	const op = "storage.postgres.Assignment"

	query := `
		SELECT *
		FROM assignments
		WHERE id = $1
	`

	var assignment models.Assignment

	err := s.db.QueryRowContext(ctx, query, assignmentID).Scan(
		&assignment.ID,
		&assignment.TeacherID,
		&assignment.Title,
		&assignment.Content,
		&assignment.CreatedAt,
		&assignment.UpdatedAt,
		&assignment.Deadline,
	)
	if err != nil {
		return models.Assignment{}, err
	}

	return assignment, nil
}
