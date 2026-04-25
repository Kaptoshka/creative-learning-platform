package role

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	pool *pgxpool.Pool
}

// New creates a new instance of RoleRepo.
// That used to interact with the role and permission related tables.
func New(pool *pgxpool.Pool) RoleRepo {
	return RoleRepo{pool: pool}
}

// LinkUserRole links a user to a specific role.
func (r *RoleRepo) LinkUserRole(
	ctx context.Context,
	userID int64,
	roleID int64,
) error {
	const op = "adapters.driven.postgres.role.LinkUserRole"

	query := `
		INSERT INTO user_roles
		(user_id, role_id)
		VALUES (@user_id, @role_id)`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"role_id": roleID,
	})
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

// RoleID returns the ID of the role with the given name.
func (r *RoleRepo) RoleID(
	ctx context.Context,
	role string,
) (int64, error) {
	const op = "adapters.driven.postgres.role.RoleID"

	query := `
		SELECT id
		FROM roles
		WHERE role = @role`

	var roleID int64

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"role": role,
	}).Scan(&roleID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf(
				"%s: %v", op, domain.ErrRoleNotFound,
			)
		}

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	return roleID, nil
}

// UserRole returns the role of the user with the given ID.
func (r *RoleRepo) UserRole(
	ctx context.Context,
	userID int64,
) (string, error) {
	const op = "adapters.driven.postgres.role.UserRole"

	query := `
		SELECT r.role
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = @user_id`

	var role string

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf(
				"%s: %v", op, domain.ErrRoleNotFound,
			)
		}
		return "", fmt.Errorf("%s: %v", op, err)
	}

	return role, nil
}

// Scope returns the permission scope of the user with the given ID.
func (r *RoleRepo) Scope(
	ctx context.Context,
	userID int64,
) ([]string, error) {
	const op = "adapters.driven.postgres.role.Scope"

	query := `
		SELECT p.slug
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = @user_id`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	scope, err := pgx.CollectRows(rows,
		func(row pgx.CollectableRow) (string, error) {
			var slug string
			return slug, row.Scan(&slug)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return scope, nil
}
