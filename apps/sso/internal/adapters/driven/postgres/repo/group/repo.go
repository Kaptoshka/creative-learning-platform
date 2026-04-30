package group

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) GroupRepo {
	return GroupRepo{pool: pool}
}

func (r *GroupRepo) UserGroups(
	ctx context.Context,
	userID uuid.UUID,
) ([]uuid.UUID, error) {
	const op = "adapters.driven.postgres.group.UserGroups"

	const query = `
		SELECT group_id
		FROM user_groups
		WHERE user_id = @user_id`

	rows, err := r.pool.Query(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var groupIDs []uuid.UUID

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		groupIDs = append(groupIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return groupIDs, nil
}
