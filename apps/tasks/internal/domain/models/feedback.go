package models

import (
	"time"

	"tasks/internal/domain"

	"github.com/google/uuid"
)

type Feedback struct {
	ID          uuid.UUID `db:"id"`
	VersionID   uuid.UUID `db:"version_id"`
	GraderID    uuid.UUID `db:"grader_id"`
	TextContent *string   `db:"text_content"`
	Payload     domain.JSONB     `db:"payload"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type FeedbackDTO struct {
	VersionID    uuid.UUID `db:"version_id"`
	SubmissionID uuid.UUID `db:"submission_id"`
	TextContent  string    `db:"text_content"`
	Payload      domain.JSONB     `db:"payload"`
	IsPublished  bool      `db:"is_published"`
}
