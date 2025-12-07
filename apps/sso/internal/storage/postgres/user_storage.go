package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/domain/models"
	"sso/internal/storage"

	pgConn "github.com/jackc/pgx/v5/pgconn"
)

type PostgresUserStorage struct {
	db *sql.DB
}

// NewUserStorage creates a new instance of PostgresUserStorage.
// That used to interact with the user table.
func NewUserStorage(db *sql.DB) *PostgresUserStorage {
	return &PostgresUserStorage{
		db: db,
	}
}

// SaveUser saves a new user to the database.
func (s *PostgresUserStorage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
	firstName string,
	lastName string,
	middleName string,
) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64

	query := `
		INSERT INTO users
		(email, pass_hash, first_name, last_name, middle_name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`
	err := s.db.QueryRowContext(ctx, query, email, passHash, firstName, lastName, middleName).Scan(&id)
	if err != nil {
		var pgErr *pgConn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %v", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return id, nil
}

// User returns the user with the given email.
func (s *PostgresUserStorage) User(
	ctx context.Context,
	email string,
) (models.User, error) {
	const op = "storage.postgres.User"

	query := `
		SELECT id, email, pass_hash, first_name, last_name, middle_name
		FROM users
		WHERE email = $1
	`

	res := s.db.QueryRowContext(ctx, query, email)
	var user models.User
	err := res.Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %v", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %v", op, err)
	}

	return user, nil
}
