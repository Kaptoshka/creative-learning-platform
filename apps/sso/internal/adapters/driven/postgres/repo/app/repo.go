package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/jackc/pgx/v5"
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
	appID int,
) (models.App, error) {
	const op = "adapters.driven.postgres.app.App"

	query := `
		SELECT
			id,
			name,
			secret
		FROM apps
		WHERE id = @id`

	var app models.App

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": appID,
	}).Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf(
				"%s: %v", op, domain.ErrAppNotFound,
			)
		}

		return models.App{}, fmt.Errorf("%s: %v", op, err)
	}

	return app, nil
}
