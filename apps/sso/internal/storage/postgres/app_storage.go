package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/domain/models"
	"sso/internal/storage"
)

type PostgresAppStorage struct {
	db *sql.DB
}

func NewAppStorage(db *sql.DB) *PostgresAppStorage {
	return &PostgresAppStorage{
		db: db,
	}
}

// App returns the app with the given ID.
func (s *PostgresAppStorage) App(
	ctx context.Context,
	appID int,
) (models.App, error) {
	const op = "storage.postgres.App"

	query := `
		SELECT *
		FROM apps
		WHERE id = $1
	`

	res := s.db.QueryRowContext(ctx, query, appID)

	var app models.App

	err := res.Scan(
		&app.ID,
		&app.Name,
		&app.Secret,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %v", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %v", op, err)
	}

	return app, nil
}
