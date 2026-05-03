package role

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) Repo {
	return Repo{pool: pool}
}

func (r *Repo) LinkUserRole(
	ctx context.Context,
	userID uuid.UUID,
	roleID uuid.UUID,
) error {
	const op = "adapters.driven.postgres.role.LinkUserRole"

	const query = `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (@user_id, @role_id)`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id": userID,
		"role_id": roleID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) RoleID(
	ctx context.Context,
	role string,
) (uuid.UUID, error) {
	const op = "adapters.driven.postgres.role.RoleID"

	const query = `
		SELECT id
		FROM roles
		WHERE role = @role`

	var roleID uuid.UUID

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"role": role,
	}).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return roleID, nil
}

func (r *Repo) UserRole(
	ctx context.Context,
	userID uuid.UUID, // ← было int64
) (string, error) {
	const op = "adapters.driven.postgres.role.UserRole"

	const query = `
		SELECT r.role
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = @user_id`

	var role string

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	}).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, domain.ErrRoleNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}

func (r *Repo) Scope(
	ctx context.Context,
	userID uuid.UUID,
) ([]string, error) {
	const op = "adapters.driven.postgres.role.Scope"

	const query = `
		SELECT p.slug
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = @user_id`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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
