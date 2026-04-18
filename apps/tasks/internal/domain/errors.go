package domain

import (
	"errors"
)

var (
	ErrTitleRequired    = errors.New("title is required")
	ErrWidgetIDRequired = errors.New("widget_id is required")
	ErrDueDateTooClose  = errors.New("due_date must be at least cutoff duration in the future")
	ErrTargetsRequired  = errors.New("at least one target is required")
	ErrTargetEmpty      = errors.New("each target must have group_id or student_id")

	ErrForbidden = errors.New("forbidden: caller is not the owner")
	ErrNotFound  = errors.New("entity not found")

	ErrNoUpdates      = errors.New("no fields to update")
	ErrInvalidDueDate = errors.New("due_date has invalid type")
	ErrForbiddenField = errors.New("field cannot be updated directly")

	ErrInvalidPageToken = errors.New("invalid or malformed page token")

	ErrInvalidStatusFilter = errors.New("invalid submission status filter")
)
