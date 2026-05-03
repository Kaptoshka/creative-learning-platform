package postgres

import (
	"errors"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapPgError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case "23505": // unique_violation
		return domain.ErrAlreadyExists
	case "23503": // foreign_key_violation
		return domain.ErrNotFound
	default:
		return err
	}
}
