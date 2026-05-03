package domain

import "errors"

var (
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrInternal         = errors.New("internal error")

	ErrSubmissionClosed       = errors.New("submission is not in progress")
	ErrSubmissionNotSubmitted = errors.New("submission is not submitted")
	ErrNoUpdates              = errors.New("no updates provided")
	ErrInvalidPageToken       = errors.New("invalid page token")
)
