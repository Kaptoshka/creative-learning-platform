package postgres

import (
	"context"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres/repo/app"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres/repo/group"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres/repo/refresh"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres/repo/role"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres/repo/user"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	app.AppRepo
	role.RoleRepo
	user.UserRepo
	group.GroupRepo
	refresh.RefreshRepo
}

// New creates a new instance of PostgreSQL storage
func New(connString string) (*Storage, error) {
	const op = "storage.postgres.New"

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return &Storage{
		pool:        pool,
		AppRepo:     app.New(pool),
		RoleRepo:    role.New(pool),
		UserRepo:    user.New(pool),
		GroupRepo:   group.New(pool),
		RefreshRepo: refresh.New(pool),
	}, nil
}

// Close closes the database connection.
func (s *Storage) Close() error {
	if s.pool != nil {
		s.pool.Close()
	}

	return nil
}
