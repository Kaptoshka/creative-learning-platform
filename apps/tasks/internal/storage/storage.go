package storage

import (
	"errors"
)

var (
	ErrAssignmentAlreadyExists = errors.New("assignment already exists")
	ErrAssignmentNotFound      = errors.New("assignment not found")
	ErrAssignmentUpdateFailed  = errors.New("assignment update failed")
	ErrSubmissionNotFound      = errors.New("submission not found")
)

type SubmissionStorage interface {
}

type AssignmentStorage interface {
}
