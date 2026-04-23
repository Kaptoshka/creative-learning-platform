package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SaveVersion struct {
	SubmissionID uuid.UUID
	Payload      json.RawMessage
	TimeSpent    time.Duration
}
