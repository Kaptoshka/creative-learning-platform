package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppRepo struct {
	pool *pgxpool.Pool
}

// New creates a new instance of AppRepo.
// That used to interact with the app table.
func New(pool *pgxpool.Pool) AppRepo {
	return AppRepo{pool: pool}
}

// App returns the app with the given ID.
func (r *AppRepo) App(
	ctx context.Context,
	appID uuid.UUID,
) (models.App, error) {
	const op = "adapters.driven.postgres.app.App"

	const query = `
		SELECT
			id,
			name,
			secret,
			description,
			is_active
		FROM apps
		WHERE id = @id
		  AND is_active = true`

	var app models.App

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": appID,
	}).Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
		&app.Description,
		&app.IsActive,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, domain.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (r *AppRepo) DeactivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "adapters.driven.postgres.app.DeactivateApp"

	const query = `
		UPDATE apps
		SET
			is_active  = false,
			updated_at = NOW()
		WHERE id = @id`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{"id": appID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AppRepo) ActivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "adapters.driven.postgres.app.ActivateApp"

	const query = `
		UPDATE apps
		SET
			is_active  = true,
			updated_at = NOW()
		WHERE id = @id`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{"id": appID})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *AppRepo) RegisterApp(
	ctx context.Context,
	name string,
	secret string,
	description string,
) (uuid.UUID, error) {
	const op = "adapters.driven.postgres.app.RegisterApp"

	const query = `
		INSERT INTO apps (name, secret, description)
		VALUES (@name, @secret, @description)
		RETURNING id`

	var id uuid.UUID

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"name":        name,
		"secret":      secret,
		"description": description,
	}).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrAppExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
