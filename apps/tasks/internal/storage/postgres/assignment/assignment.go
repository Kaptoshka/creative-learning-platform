package assignment

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

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
	var contentBytes []byte

	err := s.db.QueryRowContext(ctx, query, assignmentID).Scan(
		&assignment.ID,
		&assignment.TeacherID,
		&assignment.Title,
		&contentBytes,
		&assignment.CreatedAt,
		&assignment.UpdatedAt,
		&assignment.Deadline,
	)
	if err != nil {
		return models.Assignment{}, err
	}

	if err := json.Unmarshal(contentBytes, &assignment.Content); err != nil {
		return models.Assignment{}, err
	}

	return assignment, nil
}

func (s *PostgresAssignmentStorage) SaveAssignment(
	ctx context.Context,
	assignment models.Assignment,
) (int64, error) {
	const op = "storage.postgres.SaveASsignment"

	query := `
		INSERT INTO assignments
		(teacher_id, title, content, created_at, updated_at, deadline)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	res := s.db.QueryRowContext(
		ctx,
		query,
		assignment.TeacherID,
		assignment.Title,
		assignment.Content,
		time.Now().Unix(),
		time.Now().Unix(),
		assignment.Deadline.Unix(),
	)

	var id int64
	err := res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
