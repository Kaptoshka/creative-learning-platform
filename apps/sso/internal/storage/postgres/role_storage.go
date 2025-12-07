package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/storage"
)

type PostgresRoleStorage struct {
	db *sql.DB
}

func NewRoleStorage(db *sql.DB) *PostgresRoleStorage {
	return &PostgresRoleStorage{
		db: db,
	}
}

// LinkUserRole links a user to a specific role.
func (s *PostgresRoleStorage) LinkUserRole(
	ctx context.Context,
	userID int64,
	roleID int64,
) error {
	const op = "storage.postgres.SaveUserRole"

	query := `
		INSERT INTO user_roles
		(user_id, role_id)
		VALUES ($1, $2)
	`

	_, err := s.db.ExecContext(ctx, query, userID, roleID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

// RoleID returns the ID of the role with the given name.
func (s *PostgresRoleStorage) RoleID(
	ctx context.Context,
	role string,
) (int64, error) {
	const op = "storage.postgres.RoleID"

	query := `
		SELECT id
		FROM roles
		WHERE role = $1
	`

	res := s.db.QueryRowContext(ctx, query, role)

	var roleID int64
	err := res.Scan(&roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("%s: %v", op, storage.ErrRoleNotFound)
		}

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return roleID, nil
}

// UserRole returns the role of the user with the given ID.
func (s *PostgresRoleStorage) UserRole(
	ctx context.Context,
	userID int64,
) (string, error) {
	const op = "storage.postgres.UserRole"

	query := `
		SELECT r.role
		FROM roles r
		INNER JOIN user_roles ur
		ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	res := s.db.QueryRowContext(ctx, query, userID)

	var role string
	err := res.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %v", op, storage.ErrRoleNotFound)
		}

		return "", fmt.Errorf("%s: %v", op, err)
	}

	return role, nil
}

// Scope returns the permission scope of the user with the given ID.
func (s *PostgresRoleStorage) Scope(
	ctx context.Context,
	userID int64,
) ([]string, error) {
	const op = "storage.postgres.Scope"

	query := `
		SELECT slug
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	var scope []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
		scope = append(scope, slug)
	}

	return scope, nil
}
