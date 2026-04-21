package domain

import (
	"fmt"
	"time"
)

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
	MaxPageSizeLimit     = 100
)

const (
	CutoffDuration = 5 * time.Minute
)

var validSubmissionStatuses = map[SubmissionStatus]struct{}{
	StatusNotSpecified: {},
	StatusNotStarted:   {},
	StatusInProgress:   {},
	StatusSubmitted:    {},
	StatusGraded:       {},
	StatusReturned:     {},
}

func ValidateSubmissionStatus(s SubmissionStatus) error {
	if _, ok := validSubmissionStatuses[s]; !ok {
		return fmt.Errorf("%w: %q", ErrInvalidStatusFilter, s)
	}
	return nil
}

var allowedTemplateFields = map[string]struct{}{
	"title":         {},
	"description":   {},
	"widget_id":     {},
	"widget_config": {},
	"due_date":      {},
	"updated_at":    {},
}

func ValidateTemplateField(f string) error {
	if _, allowed := allowedTemplateFields[f]; !allowed {
		return fmt.Errorf("forbidden field %q", f)
	}
	return nil
}
