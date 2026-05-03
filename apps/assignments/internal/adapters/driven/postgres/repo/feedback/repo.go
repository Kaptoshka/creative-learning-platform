package feedback

import (
	"context"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedbackRepo struct {
	pool *pgxpool.Pool
}

// New creates a new FeedbackRepo instance.
// That used to interact with the feedbacks table.
func New(pool *pgxpool.Pool) FeedbackRepo {
	return FeedbackRepo{pool: pool}
}

func (r *FeedbackRepo) CreateFeedback(
	ctx context.Context,
	feedback models.Feedback,
	changeSubmissionStatus *domain.SubmissionStatus,
) error {
	const op = "storage.postgres.feedbacks.CreateFeedback"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	const feedbackQuery = `
		INSERT INTO feedbacks (
			id,
			version_id,
			grader_id,
			text_content,
			payload,
			is_published,
			created_at,
			updated_at
		) VALUES (
			@id,
			@version_id,
			@grader_id,
			@text_content,
			@payload,
			@is_published,
			@created_at,
			@updated_at
		)
	`

	_, err = tx.Exec(ctx, feedbackQuery, pgx.NamedArgs{
		"id":           feedback.ID,
		"version_id":   feedback.VersionID,
		"grader_id":    feedback.GraderID,
		"text_content": feedback.TextContent,
		"payload":      feedback.Payload,
		"is_published": feedback.IsPublished,
		"created_at":   feedback.CreatedAt,
		"updated_at":   feedback.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("%s: insert feedback: %w", op, err)
	}

	if changeSubmissionStatus != nil {
		const updateQuery = `
			UPDATE submissions
			SET status = @status
			WHERE id IN (
				SELECT sv.submission_id
				FROM submission_versions sv
				WHERE sv.id = @version_id
			)`

		_, err = tx.Exec(ctx, updateQuery, pgx.NamedArgs{
			"status":     *changeSubmissionStatus,
			"version_id": feedback.VersionID,
		})
		if err != nil {
			return fmt.Errorf("%s: update submission status: %w", op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"%s: commit tx: %w",
			op,
			err,
		)
	}

	return nil
}

func (r *FeedbackRepo) GetFeedbacksBySubmission(
	ctx context.Context,
	submissionID uuid.UUID,
) ([]models.Feedback, error) {
	const op = "storage.postgres.feedbacks.GetFeedbacksBySubmission"

	const query = `
		SELECT
			f.id,
			f.version_id,
			f.grader_id,
			f.text_content,
			f.payload,
			f.is_published,
			f.created_at,
			f.updated_at
		FROM feedbacks f
		JOIN submission_versions sv ON f.version_id = sv.id
		WHERE sv.submission_id = @submission_id
		ORDER BY f.created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"submission_id": submissionID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	feedbacks, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (models.Feedback, error) {
			var f models.Feedback
			err := row.Scan(
				&f.ID,
				&f.VersionID,
				&f.GraderID,
				&f.TextContent,
				&f.Payload,
				&f.IsPublished,
				&f.CreatedAt,
				&f.UpdatedAt,
			)
			return f, err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: scan: %w", op, err)
	}

	return feedbacks, nil
}
