package models

type JSONB map[string]any

type SubmissionStatus string

const (
	StatusInProgress SubmissionStatus = "IN_PROGRESS"
	StatusSubmitted  SubmissionStatus = "SUBMITTED"
	StatusGraded     SubmissionStatus = "GRADED"
	StatusReturned   SubmissionStatus = "RETURNED"
)
