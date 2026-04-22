package assignment

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssignmentRepo struct {
	pool *pgxpool.Pool
}

// New creates a new AssignmentRepo instance.
// That used to interact with the assignments table.
func New(pool *pgxpool.Pool) *AssignmentRepo {
	return &AssignmentRepo{pool: pool}
}

func (r *AssignmentRepo) CreateAssignment(
	ctx context.Context,
	tmpl models.AssignmentTemplate,
	targets []*models.AssignmentTarget,
) error {
	const op = "storage.postgres.assignment.CreateAssignment"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	const tmplQuery = `
		INSERT INTO assignment_templates (
			id,
			creator_id,
			title,
			description,
			widget_id,
			widget_config,
			due_date,
			created_at,
			updated_at
		) VALUES (
			@id,
			@creator_id,
			@title,
			@description,
			@widget_id,
			@widget_config,
			@due_date,
			@created_at,
			@updated_at
		)
	`

	tmplArgs := pgx.NamedArgs{
		"id":            tmpl.ID,
		"creator_id":    tmpl.CreatorID,
		"title":         tmpl.Title,
		"description":   tmpl.Description,
		"widget_id":     tmpl.WidgetID,
		"widget_config": tmpl.WidgetConfig,
		"due_date":      tmpl.DueDate,
		"created_at":    tmpl.CreatedAt,
		"updated_at":    tmpl.UpdatedAt,
	}

	if _, err := tx.Exec(ctx, tmplQuery, tmplArgs); err != nil {
		return fmt.Errorf("%s: insert template: %w", op, err)
	}

	if len(targets) > 0 {
		const targetQuery = `
			INSERT INTO assignment_targets (
				id,
				template_id,
				group_id,
				student_id,
				created_at,
				updated_at
			) VALUES (
				@id,
				@template_id,
				@group_id,
				@student_id,
				@created_at,
				@updated_at
			)
		`

		batch := &pgx.Batch{}
		for _, t := range targets {
			batch.Queue(targetQuery, pgx.NamedArgs{
				"id":          t.ID,
				"template_id": t.TemplateID,
				"group_id":    t.GroupID,
				"student_id":  t.StudentID,
				"created_at":  t.CreatedAt,
				"updated_at":  t.UpdatedAt,
			})
		}

		br := tx.SendBatch(ctx, batch)
		defer br.Close()

		for range targets {
			if _, err := br.Exec(); err != nil {
				return fmt.Errorf("%s: batch insert target: %w", op, err)
			}
		}

		if err := br.Close(); err != nil {
			return fmt.Errorf("%s: close batch: %w", op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit tx: %w", op, err)
	}

	return nil
}

func (r *AssignmentRepo) UpdateAssignment(
	ctx context.Context,
	id uuid.UUID,
	updates map[string]any,
	newTargets []models.AssignmentTarget,
) (*models.AssignmentTemplate, error) {
	const op = "storage.postgres.assignments.UpdateAssignment"

	setClauses := make([]string, 0, len(updates))
	args := pgx.NamedArgs{
		"id": id,
	}

	for field, value := range updates {
		if err := domain.ValidateTemplateField(field); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = @%s", field, field))
		args[field] = value
	}

	sort.Strings(setClauses)

	updateQuery := fmt.Sprintf(`
		UPDATE assignment_templates
		SET %s
		WHERE id = @id
		RETURNING
		id,
		creator_id,
		title,
		description,
		widget_id,
		widget_config,
		due_date,
		created_at,
		updated_at`,
		strings.Join(setClauses, ", "),
	)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, updateQuery, args)

	var tmpl *models.AssignmentTemplate
	var widgetConfig []byte

	err = row.Scan(
		&tmpl.ID,
		&tmpl.CreatorID,
		&tmpl.Title,
		&tmpl.Description,
		&tmpl.WidgetID,
		&widgetConfig,
		&tmpl.DueDate,
		&tmpl.CreatedAt,
		&tmpl.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: scan template: %w", op, err)
	}

	tmpl.WidgetConfig = widgetConfig

	if len(newTargets) > 0 {
		const deleteQuery = `
			DELETE FROM assignment_targets
			WHERE template_id = @template_id
		`

		if _, err := tx.Exec(ctx, deleteQuery, pgx.NamedArgs{
			"template_id": id,
		}); err != nil {
			return nil, fmt.Errorf("%s: delete targets: %w", op, err)
		}

		const targetQuery = `
			INSERT INTO assignment_targets (
				id,
				template_id,
				group_id,
				student_id,
				created_at,
				updated_at
			) VALUES (
				@id,
				@template_id,
				@group_id,
				@student_id,
				@created_at,
				@updated_at
			)
		`
		batch := &pgx.Batch{}
		for _, t := range newTargets {
			batch.Queue(targetQuery, pgx.NamedArgs{
				"id":          t.ID,
				"template_id": t.TemplateID,
				"group_id":    t.GroupID,
				"student_id":  t.StudentID,
				"created_at":  t.CreatedAt,
				"updated_at":  t.UpdatedAt,
			})
		}
		br := tx.SendBatch(ctx, batch)
		defer br.Close()

		for range newTargets {
			if _, err := br.Exec(); err != nil {
				return nil, fmt.Errorf("%s: batch insert target: %w", op, err)
			}
		}

		if err := br.Close(); err != nil {
			return nil, fmt.Errorf("%s: close batch: %w", op, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("%s: commit tx: %w", op, err)
	}

	return tmpl, nil
}

func (r *AssignmentRepo) DeleteAssignment(
	ctx context.Context,
	id uuid.UUID,
) error {
	const op = "storage.postgres.assignments.DeleteAssignment"

	const query = `
		DELETE FROM assignment_templates
		WHERE id = @id`

	result, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, domain.ErrNotFound)
	}

	return nil
}

func (r *AssignmentRepo) GetAssignmentByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.AssignmentTemplate, []*models.AssignmentTarget, error) {
	const op = "storage.postgres.assignments.GetAssignmentByID"

	const tmplQuery = `
			SELECT
				id,
				creator_id,
				title,
				description,
				widget_id,
				widget_config,
				due_date,
				created_at,
				updated_at
			FROM assignment_templates
			WHERE id = @id`

	var tmpl models.AssignmentTemplate

	err := r.pool.QueryRow(ctx, tmplQuery, pgx.NamedArgs{
		"id": id,
	}).Scan(
		&tmpl.ID,
		&tmpl.CreatorID,
		&tmpl.Title,
		&tmpl.Description,
		&tmpl.WidgetID,
		&tmpl.WidgetConfig,
		&tmpl.DueDate,
		&tmpl.CreatedAt,
		&tmpl.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	const targetsQuery = `
			SELECT
				id,
				template_id,
				group_id,
				student_id,
				created_at,
				updated_at
			FROM assignment_targets
			WHERE template_id = @template_id`

	rows, err := r.pool.Query(ctx, targetsQuery, pgx.NamedArgs{
		"template_id": id,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("%s: query targets: %w", op, err)
	}
	defer rows.Close()

	targets, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (*models.AssignmentTarget, error) {
			var t *models.AssignmentTarget
			err := row.Scan(
				&t.ID,
				&t.TemplateID,
				&t.GroupID,
				&t.StudentID,
				&t.CreatedAt,
				&t.UpdatedAt,
			)
			return t, err
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: scan targets: %w", op, err)
	}

	return &tmpl, targets, nil
}

func (r *AssignmentRepo) ListTemplatesByCreator(
	ctx context.Context,
	creatorID uuid.UUID,
	limit int,
	offset int,
) ([]*models.AssignmentTemplate, error) {
	const op = "storage.postgres.assignments.ListTemplatesByCreator"

	const query = `
		SELECT
			id,
			creator_id,
			title,
			description,
			widget_id,
			widget_config,
			due_date,
			created_at,
			updated_at
		FROM assignment_templates
		WHERE creator_id = @creator_id
		ORDER BY id DESC
		LIMIT @limit
		OFFSET @offset`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"creator_id": creatorID,
		"limit":      limit,
		"offset":     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	templates, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (*models.AssignmentTemplate, error) {
			var t *models.AssignmentTemplate
			err := row.Scan(
				&t.ID,
				&t.CreatorID,
				&t.Title,
				&t.Description,
				&t.WidgetID,
				&t.WidgetConfig,
				&t.DueDate,
				&t.CreatedAt,
				&t.UpdatedAt,
			)
			return t, err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: scan: %w", op, err)
	}

	return templates, nil
}

func (r *AssignmentRepo) CreateSubmission(
	ctx context.Context,
	sub models.Submission,
) error {
	const op = "storage.postgres.assignments.CreateSubmission"

	const query = `
		INSERT INTO submissions (
			id,
			template_id,
			student_id,
			status,
			started_at
		) VALUES (
			@id,
			@template_id,
			@student_id,
			@status,
			@started_at
		)`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"id":          sub.ID,
		"template_id": sub.TemplateID,
		"student_id":  sub.StudentID,
		"status":      sub.Status,
		"started_at":  sub.StartedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AssignmentRepo) GetSubmissionByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Submission, error) {
	const op = "storage.postgres.assignments.GetSubmissionByID"

	const submissionQuery = `
		SELECT
			id,
			template_id,
			student_id,
			status,
			started_at,
			submitted_at
		FROM submissions
		WHERE id = @id
	`

	var submission models.Submission

	err := r.pool.QueryRow(ctx, submissionQuery, pgx.NamedArgs{
		"id": id,
	}).Scan(
		&submission.ID,
		&submission.TemplateID,
		&submission.StudentID,
		&submission.Status,
		&submission.StartedAt,
		&submission.SubmittedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &submission, nil
}

func (r *AssignmentRepo) AddSubmissionVersion(
	ctx context.Context,
	version models.SubmissionVersion,
	updateParentStatus bool,
) error {
	const op = "storage.postgres.assignments.AddSubmissionVersion"
	const query = `
		INSERT INTO submission_versions (
			id,
			submission_id,
			version_number,
			payload,
			time_stent_seconds,
			is_autosave,
			created_at,
			updated_at,
		) VALUES (
			@id,
			@submission_id,
			@version_number,
			@payload,
			@time_stent_seconds,
			@is_autosave,
			@created_at,
			@updated_at,
		)
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = r.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{
			"id":                 version.ID,
			"submission_id":      version.SubmissionID,
			"version_number":     version.VersionNumber,
			"payload":            version.Payload,
			"time_stent_seconds": version.TimeSpentSeconds,
			"is_autosave":        version.IsAutosave,
			"created_at":         version.CreatedAt,
			"updated_at":         version.UpdatedAt,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("%s: %w", op, domain.ErrAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	if !version.IsAutosave {
		const subQuery = `
			UPDATE submissions
				SET
					status = @status,
					submitted_at = @submitted_at
				WHERE id = @submission_id
		`

		_, err = r.pool.Exec(ctx, subQuery, pgx.NamedArgs{
			"status":       domain.StatusSubmitted,
			"submitted_at": time.Now().UTC(),
		})
		if err != nil {
			return fmt.Errorf(
				"%s: failed to update submission: %w",
				op,
				err,
			)
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

func (r *AssignmentRepo) ListSubmissionByTemplate(
	ctx context.Context,
	templateID uuid.UUID,
	limit int,
	offset int,
	filter domain.SubmissionStatus,
) ([]models.SubmissionVersion, error) {
	const op = "storage.postgres.assignments.ListSubmissionByTemplate"

	const query = `

	`
}

func (r *AssignmentRepo) GetFeedbacksBySubmission(
	ctx context.Context,
	submissionID uuid.UUID,
) ([]*models.Feedback, error) {
	const op = "storage.postgres.assignments.GetFeedbacksBySubmission"

	const query = `
		SELECT
			id,
			version_id,
			grader_id,
			text_content,
			payload,
			is_published,
			created_at,
			updated_at
		FROM feedbacks
		WHERE version_id IN (
			SELECT id
				FROM submission_versions
				WHERE submission_id = @submission_id
		)
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
		func(row pgx.CollectableRow) (*models.Feedback, error) {
			var f *models.Feedback
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

func (r *AssignmentRepo) CreateFeedback(
	ctx context.Context,
	feedback models.Feedback,
	changeSubmissionStatus *domain.SubmissionStatus,
) error {
	const op = "storage.postgres.assignments.CreateFeedback"

	const feedbackQuery = `
		INSERT INTO feedbacks (
			id,
			version_id,
			grader_id,
			text_content,
			payload,
			is_published,
			created_at,
			updated_at,
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

	const submissionQuery = `
		UPDATE submissions
			SET status = @status
			WHERE id IN (
				SELECT sv.submission_id
					FROM submission_versions sv
					WHERE sv.id = @version_id
			)
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		feedbackQuery,
		pgx.NamedArgs{
			"id":           feedback.ID,
			"version_id":   feedback.VersionID,
			"grader_id":    feedback.GraderID,
			"text_content": feedback.TextContent,
			"payload":      feedback.Payload,
			"is_published": feedback.IsPublished,
			"created_at":   feedback.CreatedAt,
			"updated_at":   feedback.UpdatedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: insert feedback: %w", op, err)
	}

	_, err = tx.Exec(
		ctx,
		submissionQuery,
		pgx.NamedArgs{
			"status":     changeSubmissionStatus,
			"version_id": feedback.VersionID,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: update submission: %w", op, err)
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
