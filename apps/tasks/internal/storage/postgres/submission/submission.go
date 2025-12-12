package submission

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"tasks/internal/domain/models"
	"tasks/internal/storage"
)

type PostgresSubmissionStorage struct {
	db *sql.DB
}

// New creates a new PostgresSubmissionStorage instance.
// That used to interact with the submissions table.
func New(db *sql.DB) *PostgresSubmissionStorage {
	return &PostgresSubmissionStorage{db: db}
}

func (s *PostgresSubmissionStorage) SubmissionByAssignmentID(
	ctx context.Context,
	assignmentID int64,
) ([]models.Submission, error) {
	const op = "storage.postgres.SubmissionByAssignmentID"

	query := `
		SELECT *
		FROM submissions
		WHERE assignment_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, assignmentID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var submissions []models.Submission
	var contentBytes []byte
	for rows.Next() {
		var submission models.Submission
		if err := rows.Scan(
			&submission.ID,
			&submission.AssignmentID,
			&submission.StudentID,
			&contentBytes,
			&submission.StartedAt,
			&submission.SubmittedAt,
			&submission.Status,
			&submission.Feedback,
		); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		if err := json.Unmarshal(contentBytes, &submission.Content); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}

		submissions = append(submissions, submission)
	}

	return submissions, nil
}

func (s *PostgresSubmissionStorage) SubmissionByAssignmentIDAndStudentID(
	ctx context.Context,
	assignmentID int64,
	studentID int64,
) (models.Submission, error) {
	const op = "storage.postgres.SubmissionByAssignmentIDAndStudentID"

	query := `
		SELECT *
		FROM submissions
		WHERE assignment_id = $1 AND student_id = $2
	`

	res := s.db.QueryRowContext(ctx, query, assignmentID, studentID)

	var submission models.Submission
	var contentBytes []byte
	err := res.Scan(
		&submission.ID,
		&submission.AssignmentID,
		&submission.StudentID,
		&contentBytes,
		&submission.StartedAt,
		&submission.SubmittedAt,
		&submission.Status,
		&submission.Feedback,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Submission{}, fmt.Errorf("%s: %v", op, storage.ErrSubmissionNotFound)
		}
		return models.Submission{}, fmt.Errorf("%s: %v", op, err)
	}

	if err := json.Unmarshal(contentBytes, &submission.Content); err != nil {
		return models.Submission{}, fmt.Errorf("%s: %v", op, err)
	}

	return submission, nil
}
