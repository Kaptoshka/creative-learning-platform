package dto

import (
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type StudentItem struct {
	AssignmentID uuid.UUID
	Title        string
	WidgetType   string
	DueDate      *time.Time
	CreatedAt    time.Time
	SubmissionID *uuid.UUID
	Status       *domain.SubmissionStatus
	StartedAt    *time.Time
	SubmittedAt  *time.Time
	HasFeedback  bool
}
