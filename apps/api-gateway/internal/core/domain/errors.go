package domain

import "errors"

var (
	ErrUnauthenticated  = errors.New("unauthenticated")
	ErrPermissionDenied = errors.New("permission denied")
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvalidArgument  = errors.New("invalid argument")
	ErrInternal         = errors.New("internal error")
)
