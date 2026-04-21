package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Feedback struct {
	VersionID    uuid.UUID       `db:"version_id"`
	SubmissionID uuid.UUID       `db:"submission_id"`
	TextContent  string          `db:"text_content"`
	Payload      json.RawMessage `db:"payload"`
	IsPublished  bool            `db:"is_published"`
}
