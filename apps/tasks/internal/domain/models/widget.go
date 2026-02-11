package models

import (
	"github.com/google/uuid"
)

type Widget struct {
	ID               uuid.UUID `db:"id"`
	Type             string    `db:"type"`
	Version          int       `db:"version"`
	ConfigSchema     JSONB     `db:"config_schema"`
	SubmissionSchema JSONB     `db:"submission_schema"`
}
