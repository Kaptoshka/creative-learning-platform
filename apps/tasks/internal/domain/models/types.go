package models

type JSONB map[string]any

type SubmissionStatus string

const (
	StatusNotSpecified SubmissionStatus = "NOT_SPECIFIED"
	StatusNotStarted   SubmissionStatus = "NOT_STARTED"
	StatusInProgress   SubmissionStatus = "IN_PROGRESS"
	StatusSubmitted    SubmissionStatus = "SUBMITTED"
	StatusGraded       SubmissionStatus = "GRADED"
	StatusReturned     SubmissionStatus = "RETURNED"
)

const (
	DefaultPageSizeLimit = 10
)
