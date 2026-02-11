package models

import (
	"time"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `db:"id"`
	VersionID   uuid.UUID `db:"version_id"`
	GraderID    uuid.UUID `db:"grader_id"`
	TextContent *string   `db:"text_content"`
	Payload     JSONB     `db:"payload"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
