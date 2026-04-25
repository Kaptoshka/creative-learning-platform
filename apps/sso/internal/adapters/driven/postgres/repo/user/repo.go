package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

// New creates a new instance of UserRepo.
// That used to interact with the user table.
func New(pool *pgxpool.Pool) UserRepo {
	return UserRepo{pool: pool}
}

// SaveUser saves a new user to the database.
func (r *UserRepo) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
	firstName string,
	lastName string,
	middleName string,
) (int64, error) {
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

	var id int64

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
			return 0, fmt.Errorf("%s: %v", op, domain.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

// User returns the user with the given email.
func (r *UserRepo) User(
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf(
				"%s: %v", op, domain.ErrUserNotFound,
			)
		}

		return models.User{}, fmt.Errorf("%s: %v", op, err)
	}

	return user, nil
}
