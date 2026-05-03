package refresh

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"
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

func (r *Repo) Save(
	ctx context.Context,
	userID uuid.UUID,
	tokenHash string,
	expiresAt time.Time,
) (uuid.UUID, error) {
	const op = "adapters.driven.postgres.refresh.Save"

	const query = `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at
		) VALUES (
			@user_id,
			@token_hash,
			@expires_at
		) RETURNING id`

	var id uuid.UUID

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"user_id":    userID,
		"token_hash": tokenHash,
		"expires_at": expiresAt,
	}).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repo) ByHash(
	ctx context.Context,
	tokenHash string,
) (models.RefreshToken, error) {
	const op = "adapters.driven.postgres.refresh.ByHash"

	const query = `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			revoked_at,
			created_at
		FROM refresh_tokens
		WHERE token_hash = @token_hash`

	var token models.RefreshToken

	err := r.pool.QueryRow(ctx, query, pgx.NamedArgs{
		"token_hash": tokenHash,
	}).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.RevokedAt,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.RefreshToken{}, fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		return models.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (r *Repo) Revoke(
	ctx context.Context,
	tokenID uuid.UUID,
) error {
	const op = "adapters.driven.postgres.refresh.Revoke"

	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE id = @id
		  AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"id": tokenID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) RevokeAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	const op = "adapters.driven.postgres.refresh.RevokeAll"

	const query = `
		UPDATE refresh_tokens
		SET revoked_at = NOW()
		WHERE user_id = @user_id
		  AND revoked_at IS NULL`

	_, err := r.pool.Exec(ctx, query, pgx.NamedArgs{
		"user_id": userID,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) DeleteExpired(ctx context.Context) error {
	const op = "adapters.driven.postgres.refresh.DeleteExpired"

	const query = `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW() - INTERVAL '30 days'`

	_, err := r.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
