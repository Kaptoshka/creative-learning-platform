package dto

import (
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"

	"github.com/google/uuid"
)

type SaveVersion struct {
	SubmissionID uuid.UUID
	Payload      domain.JSONB
	TimeSpent    time.Duration
}
