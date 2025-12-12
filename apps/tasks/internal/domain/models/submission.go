package models

import (
	"time"
)

type Submission struct {
	ID           int64
	AssignmentID int64                  `json:"assignment_id"`
	StudentID    int64                  `json:"student_id"`
	Content      map[string]interface{} `json:"content"`
	StartedAt    time.Time              `json:"started_at"`
	SubmittedAt  time.Time              `json:"submitted_at"`
	Status       string                 `json:"status"`
	Feedback     string                 `json:"feedback"`
}
