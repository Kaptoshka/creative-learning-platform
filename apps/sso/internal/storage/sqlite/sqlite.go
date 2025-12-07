package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New creates a new instance of SQLite storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
	firstName string,
	lastName string,
	middleName string,
) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmp, err := s.db.Prepare(
		"INSERT INTO users (email, pass_hash, first_name, last_name, middle_name) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmp.ExecContext(ctx, email, passHash, firstName, lastName, middleName)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveEnrollment(ctx context.Context, userID int64, roleID int) error {
	const op = "storage.sqlite.SaveEnrollment"

	stmp, err := s.db.Prepare(
		"INSERT INTO enrollments (user_id, role_id) VALUES (?, ?)",
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmp.ExecContext(ctx, userID, roleID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// User return user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmp, err := s.db.Prepare(
		"SELECT id, email, pass_hash, first_name, last_name, middle_name FROM users WHERE email = ?",
	)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmp.QueryRowContext(ctx, email)

	var user models.User
	err = res.Scan(&user.ID, &user.Email, &user.PassHash, &user.FirstName, &user.LastName, &user.MiddleName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

// UserExists returns true if user exists
func (s *Storage) UserExists(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.UserExists"

	stmp, err := s.db.Prepare(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)",
	)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmp.QueryRowContext(ctx, userID)

	var exists bool
	err = row.Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (s *Storage) ListUsers(ctx context.Context, roleQuery string, searchQuery string) ([]models.User, error) {
	const op = "storage.sqlite.ListUsers"

	query := `
	SELECT u.id, u.email, u.first_name, u.last_name, u.middle_name
	FROM users u
	JOIN enrollments e ON u.id = e.user_id
	JOIN roles r ON e.role_id = r.id
	WHERE r.role = ? AND (u.first_name LIKE ? OR u.last_name LIKE ? OR u.email LIKE ?)
	`

	stmp, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	searchPattern := "%" + searchQuery + "%"
	rows, err := stmp.QueryContext(ctx, roleQuery, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.MiddleName); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}

// UserRole returns role of the user
func (s *Storage) UserRole(ctx context.Context, userID int64) (string, error) {
	const op = "storage.sqlite.UserRole"

	stmp, err := s.db.Prepare(
		"SELECT r.role FROM roles r INNER JOIN enrollments en ON r.id = en.role_id WHERE en.user_id = ?",
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res := stmp.QueryRowContext(ctx, userID)

	var role string
	err = res.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmp, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	res := stmp.QueryRowContext(ctx, appID)

	var app models.App

	err = res.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
