package models

import (
	"encoding/json"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type Submission struct {
	ID          uuid.UUID               `db:"id"`
	TemplateID  uuid.UUID               `db:"template_id"`
	StudentID   uuid.UUID               `db:"student_id"`
	Status      domain.SubmissionStatus `db:"status"`
	StartedAt   time.Time               `db:"started_at"`
	SubmittedAt *time.Time              `db:"submitted_at"`
	LastVersion *SubmissionVersionLight
}

type SubmissionVersion struct {
	ID               uuid.UUID       `db:"id"`
	SubmissionID     uuid.UUID       `db:"submission_id"`
	VersionNumber    int32           `db:"version_number"`
	Payload          json.RawMessage `db:"payload"`
	TimeSpentSeconds time.Duration   `db:"time_spent_seconds"`
	IsAutosave       bool            `db:"is_autosave"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}

type SubmissionVersionLight struct {
	ID            *uuid.UUID
	VersionNumber *int32
	Payload       json.RawMessage
	CreatedAt     *time.Time
}
