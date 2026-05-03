package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new instance of Repo.
// That used to interact with the user table.
func New(pool *pgxpool.Pool) Repo {
	return Repo{pool: pool}
}

// SaveUser saves a new user to the database.
func (r *Repo) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
	firstName string,
	lastName string,
	middleName string,
) (uuid.UUID, error) {
	const op = "adapters.driven.postgres.user.SaveUser"

	const query = `
		INSERT INTO users (
			email,
			pass_hash,
			first_name,
			last_name,
			middle_name
		) VALUES (
			@email,
			@pass_hash,
			@first_name,
			@last_name,
			@middle_name
		) RETURNING id`

	var id uuid.UUID

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email":       email,
		"pass_hash":   passHash,
		"first_name":  firstName,
		"last_name":   lastName,
		"middle_name": middleName,
	}).Scan(&id)
	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns the user with the given email.
func (r *Repo) User(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "adapters.driven.postgres.user.User"

	const query = `
		SELECT
			id,
			email,
			pass_hash,
			first_name,
			last_name,
			middle_name
		FROM users
		WHERE email = @email`

	var user models.User

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"email": email,
	}).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf(
				"%s: %w", op, domain.ErrUserNotFound,
			)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (r *Repo) UserByID(
	ctx context.Context,
	userID uuid.UUID,
) (models.User, error) {
	const op = "adapters.driven.postgres.user.UserByID"

	const query = `
		SELECT
			id,
			email,
			pass_hash,
			first_name,
			last_name,
			middle_name
		FROM users
		WHERE id = @id`

	var user models.User

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"id": userID,
	}).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
