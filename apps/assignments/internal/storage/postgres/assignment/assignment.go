package assignment

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"%s: begin tx: %w",
			op,
			err,
		)
	}
	defer tx.Rollback(ctx)

	const deleteQuery = `
		DELETE FROM assignment_templates
		WHERE id = @id
	`

	if _, err := tx.Exec(ctx, deleteQuery, pgx.NamedArgs{
		"id": id,
	}); err != nil {
		return fmt.Errorf(
			"%s: delete assignment: %w",
			op,
			err,
		)
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

func (r *AssignmentRepo) GetAssignmentByID(
	ctx context.Context,
	id uuid.UUID,
) (*models.AssignmentTemplate, []*models.AssignmentTarget, error) {
	const op = "storage.postgres.assignments.GetAssignmentByID"

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"%s: begin tx: %w",
			op,
			err,
		)
	}
	defer tx.Rollback(ctx)
}
