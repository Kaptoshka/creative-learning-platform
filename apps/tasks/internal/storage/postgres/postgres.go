package postgres

import (
	"context"
	"fmt"

	"tasks/internal/storage"
	"tasks/internal/storage/postgres/assignment"
	"tasks/internal/storage/postgres/submission"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
	storage.AssignmentStorage
	storage.SubmissionStorage
}

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
		pool:              pool,
		AssignmentStorage: assignment.New(pool),
		SubmissionStorage: submission.New(pool),
	}, nil
}

func (s *Storage) Close() error {
	const op = "storage.postgres.Close"

	if s.pool != nil {
		s.pool.Close()
	}

	return nil
}
