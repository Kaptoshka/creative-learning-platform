package domain

import (
	"errors"
)

var (
	// 403
	ErrForbidden = errors.New("forbidden: caller is not the owner")

	// 404
	ErrNotFound = errors.New("entity not found")

	// 409
	ErrAlreadyExists = errors.New("entity already exists")

	// template validation
	ErrTitleRequired    = errors.New("title is required")
	ErrWidgetIDRequired = errors.New("widget_id is required")
	ErrDueDateTooClose  = errors.New("due_date must be at least cutoff duration in the future")
	ErrInvalidDueDate   = errors.New("due_date has invalid type")
	ErrNoUpdates        = errors.New("no fields to update")
	ErrForbiddenField   = errors.New("field cannot be updated directly")

	// target validation
	ErrTargetsRequired = errors.New("at least one target is required")
	ErrTargetEmpty     = errors.New("each target must have group_id or student_id")

	// pagination
	ErrInvalidPageToken = errors.New("invalid or malformed page token")

	// submission status
	ErrInvalidStatusFilter = errors.New("invalid submission status filter")
)
