package submission

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionRepo struct {
	pool *pgxpool.Pool
}

// New creates a new SubmissionRepo instance.
// That used to interact with the submissions table.
func New(pool *pgxpool.Pool) *SubmissionRepo {
	return &SubmissionRepo{pool: pool}
}

func (r *SubmissionRepo) CreateSubmission(
	ctx context.Context,
	sub models.Submission,
) error {
	const op = "storage.postgres.submission.CreateSubmission"

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

func (r *SubmissionRepo) GetSubmissionByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.Submission, error) {
	const op = "storage.postgres.submission.GetSubmissionByID"

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

func (r *SubmissionRepo) AddSubmissionVersion(
	ctx context.Context,
	version models.SubmissionVersion,
	updateParentStatus bool,
) error {
	const op = "storage.postgres.submission.AddSubmissionVersion"
	const insertQuery = `
		INSERT INTO submission_versions (
			id,
			submission_id,
			version_number,
			payload,
			time_spent_seconds,
			is_autosave,
			created_at,
			updated_at
		) VALUES (
			@id,
			@submission_id,
			@version_number,
			@payload,
			@time_spent_seconds,
			@is_autosave,
			@created_at,
			@updated_at
		)
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		insertQuery,
		pgx.NamedArgs{
			"id":                 version.ID,
			"submission_id":      version.SubmissionID,
			"version_number":     version.VersionNumber,
			"payload":            version.Payload,
			"time_spent_seconds": version.TimeSpentSeconds,
			"is_autosave":        version.IsAutosave,
			"created_at":         version.CreatedAt,
			"updated_at":         version.UpdatedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("%s: insert version: %w", op, err)
	}

	if updateParentStatus {
		const updateQuery = `
			UPDATE submissions
			SET
				status = @status,
				submitted_at = @submitted_at
			WHERE id = @submission_id
		`

		_, err = tx.Exec(ctx, updateQuery, pgx.NamedArgs{
			"status":        domain.StatusSubmitted,
			"submitted_at":  time.Now().UTC(),
			"submission_id": version.SubmissionID,
		})
		if err != nil {
			return fmt.Errorf(
				"%s: update submission status: %w",
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

func (r *SubmissionRepo) GetSubmissionVersions(
	ctx context.Context,
	submissionID uuid.UUID,
) ([]models.SubmissionVersion, error) {
	const op = "storage.postgres.submission.GetSubmissionVersions"

	const query = `
		SELECT
			id,
			submission_id,
			version_number,
			payload,
			time_spent_seconds,
			is_autosave,
			created_at,
			updated_at
		FROM submission_versions
		WHERE submission_id = @submission_id
		ORDER BY version_number DESC`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"submission_id": submissionID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	versions, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (
			models.SubmissionVersion,
			error,
		) {
			var v models.SubmissionVersion
			err := row.Scan(
				&v.ID,
				&v.SubmissionID,
				&v.VersionNumber,
				&v.Payload,
				&v.TimeSpentSeconds,
				&v.IsAutosave,
				&v.CreatedAt,
				&v.UpdatedAt,
			)
			return v, err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: scan: %w", op, err)
	}

	return versions, nil
}

func (r *SubmissionRepo) ListAssignmentsForStudent(
	ctx context.Context,
	studentID uuid.UUID,
	groupID uuid.UUID,
	limit int,
	offset int,
	statusFilter domain.SubmissionStatus,
) ([]dto.StudentItem, error) {
	const op = "storage.postgres.submission.ListAssignmentsForStudent"

	const query = `
	    SELECT
	        a.id,
	        a.title,
	        w.type        AS widget_type,
	        a.due_date,
	        a.created_at,
	        s.id          AS submission_id,
	        s.status      AS submission_status,
	        s.started_at,
	        s.submitted_at,
	        EXISTS (
	            SELECT 1
	            FROM feedbacks f
	            JOIN submission_versions sv ON f.version_id = sv.id
	            WHERE sv.submission_id = s.id
	                AND f.is_published = TRUE
	        ) AS has_feedback
	    FROM assignment_templates a
	    JOIN widgets w ON w.id = a.widget_id
	    JOIN assignment_targets t ON t.template_id = a.id
	        AND (t.student_id = @student_id OR t.group_id = @group_id)
	    LEFT JOIN submissions s ON s.template_id = a.id
	        AND s.student_id = @student_id
	    WHERE (@status_filter = '' OR s.status = @status_filter)
	    ORDER BY a.id DESC
	    LIMIT  @limit
	    OFFSET @offset`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"student_id":    studentID,
		"group_id":      groupID,
		"status_filter": statusFilter,
		"limit":         limit,
		"offset":        offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (dto.StudentItem, error) {
		var item dto.StudentItem
		err := row.Scan(
			&item.AssignmentID,
			&item.Title,
			&item.WidgetType,
			&item.DueDate,
			&item.CreatedAt,
			&item.SubmissionID,
			&item.Status,
			&item.StartedAt,
			&item.SubmittedAt,
			&item.HasFeedback,
		)
		return item, err
	})
	if err != nil {
		return nil, fmt.Errorf("%s: scan: %w", op, err)
	}

	return items, nil
}

func (r *SubmissionRepo) ListSubmissionsByTemplate(
	ctx context.Context,
	templateID uuid.UUID,
	limit int,
	offset int,
	filter domain.SubmissionStatus,
) ([]models.Submission, error) {
	const op = "storage.postgres.submission.ListSubmissionsByTemplate"

	const query = `
		SELECT DISTINCT ON (s.id)
			s.id,
			s.template_id,
			s.student_id,
			s.status,
			s.started_at,
			s.submitted_at,
			sv.id             AS last_version_id,
			sv.version_number AS last_version_number,
			sv.payload        AS last_version_payload,
			sv.created_at     AS last_version_created_at
		FROM submissions s
		LEFT JOIN submission_versions sv ON sv.submission_id = s.id
		WHERE s.template_id = @template_id
			AND (@filter = '' OR s.status = @filter)
		ORDER BY s.id, sv.version_number DESC
		LIMIT @limit
		OFFSET @offset`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"template_id": templateID,
		"filter":      filter,
		"limit":       limit,
		"offset":      offset,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	submissions, err := pgx.CollectRows(
		rows,
		func(row pgx.CollectableRow) (models.Submission, error) {
			var s models.Submission
			var lastVersion models.SubmissionVersionLight
			err := row.Scan(
				&s.ID,
				&s.TemplateID,
				&s.StudentID,
				&s.Status,
				&s.StartedAt,
				&s.SubmittedAt,
				&lastVersion.ID,
				&lastVersion.VersionNumber,
				&lastVersion.Payload,
				&lastVersion.CreatedAt,
			)
			s.LastVersion = &lastVersion
			return s, err
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: scan: %w", op, err)
	}

	return submissions, nil
}
