package models

import (
	"time"

	"tasks/internal/domain"

	"github.com/google/uuid"
)

type Submission struct {
	ID          uuid.UUID        `db:"id"`
	TemplateID  uuid.UUID        `db:"template_id"`
	StudentID   uuid.UUID        `db:"student_id"`
	Status      domain.SubmissionStatus `db:"status"`
	StartedAt   time.Time        `db:"started_at"`
	SubmittedAt *time.Time       `db:"submitted_at"`
}

type SubmissionVersion struct {
	ID               uuid.UUID     `db:"id"`
	SubmissionID     uuid.UUID     `db:"submission_id"`
	VersionNumber    int32         `db:"version_number"`
	Payload          domain.JSONB         `db:"payload"`
	TimeSpentSeconds time.Duration `db:"time_spent_seconds"`
	IsAutosave       bool          `db:"is_autosave"`
	CreatedAt        time.Time     `db:"created_at"`
	UpdatedAt        time.Time     `db:"updated_at"`
}

type SubmissionItem struct {
	Submission        *Submission
	SubmissionVersion *SubmissionVersion
}
