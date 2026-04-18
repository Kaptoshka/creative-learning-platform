package models

import (
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type Widget struct {
	ID               uuid.UUID    `db:"id"`
	Type             string       `db:"type"`
	Version          int          `db:"version"`
	ConfigSchema     domain.JSONB `db:"config_schema"`
	SubmissionSchema domain.JSONB `db:"submission_schema"`
}
