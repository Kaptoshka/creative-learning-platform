package submission

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubmissionRepo struct {
	pool *pgxpool.Pool
}

// New creates a new SubmissionRepo instance.
// That used to interact with the submissions table.
func New(pool *pgxpool.Pool) *SubmissionRepo {
	return &SubmissionRepo{pool: pool}
}
