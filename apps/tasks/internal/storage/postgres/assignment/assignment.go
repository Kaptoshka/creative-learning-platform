package assignment

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"tasks/internal/domain/models"
	"tasks/internal/storage"
)

type AssignmentRepo struct {
	db *sql.DB
}

// New creates a new AssignmentRepo instance.
// That used to interact with the assignments table.
func New(db *sql.DB) *AssignmentRepo {
	return &AssignmentRepo{db: db}
}

func (r *AssignmentRepo) SaveAssignment(
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
	res := r.db.QueryRowContext(
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

func (r *AssignmentRepo) UpdateAssignment(
	ctx context.Context,
	assignmentID int64,
	updates map[string]any,
	updateTargets bool,
	studentIDs []int64,
) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if len(updates) > 0 {
		query := "UPDATE assignments SET "
		args := []any{}
		argID := 1

		i := 0
		for col, val := range updates {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("%s = %d", col, argID)
			args = append(args, val)
			argID++
			i++
		}

		query += fmt.Sprintf(" WHERE id = %d", argID)
		args = append(args, assignmentID)

		_, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return storage.ErrAssignmentUpdateFailed
		}
	}

	if updateTargets {
		_, err = tx.ExecContext(
			ctx,
			"DELETE FROM assignment_targets WHERE assignment_id = $1",
			assignmentID,
		)
		if err != nil {
			return err
		}

		for _, studentID := range studentIDs {
			_, err = tx.ExecContext(
				ctx,
				"INSERT INTO assignment_targets (assignment_id, student_id) VALUES ($1, $2)",
				assignmentID,
				studentID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *AssignmentRepo) Assignment(
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

	err := r.db.QueryRowContext(ctx, query, assignmentID).Scan(
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
