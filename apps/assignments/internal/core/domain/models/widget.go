package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Widget struct {
	ID               uuid.UUID       `db:"id"`
	Type             string          `db:"type"`
	Version          int             `db:"version"`
	ConfigSchema     json.RawMessage `db:"config_schema"`
	SubmissionSchema json.RawMessage `db:"submission_schema"`
}
