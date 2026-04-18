package domain

import (
	"fmt"
)

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

var validSubmissionStatuses = map[SubmissionStatus]struct{}{
	StatusNotSpecified: {},
	StatusNotStarted:   {},
	StatusInProgress:   {},
	StatusSubmitted:    {},
	StatusGraded:       {},
	StatusReturned:     {},
}

func validateSubmissionStatus(s SubmissionStatus) error {
	if _, ok := validSubmissionStatuses[s]; !ok {
		return fmt.Errorf("%w: %q", ErrInvalidStatusFilter, s)
	}
	return nil
}

const (
	DefaultPageSizeLimit = 10
	MaxPageSizeLimit     = 100
)
