package assignment

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssignmentRepo struct {
	pool *pgxpool.Pool
}

// New creates a new AssignmentRepo instance.
// That used to interact with the assignments table.
func New(pool *pgxpool.Pool) *AssignmentRepo {
	return &AssignmentRepo{pool: pool}
}
