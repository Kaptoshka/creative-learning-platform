package dto

import (
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type Feedback struct {
	VersionID    uuid.UUID    `db:"version_id"`
	SubmissionID uuid.UUID    `db:"submission_id"`
	TextContent  string       `db:"text_content"`
	Payload      domain.JSONB `db:"payload"`
	IsPublished  bool         `db:"is_published"`
}
