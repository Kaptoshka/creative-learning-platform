package assignment

import (
	"context"
	"log/slog"

	"tasks/internal/domain/models"
)

type AssignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
}

type AssignmentSaver interface {
	SaveAssignment() (int64, error)
}

type AssignmentProvider interface {
	AssignmentByID(ctx context.Context, id int64) (*models.Assignment, error)
	AssignmentByTeacherID(ctx context.Context, teacherID int64) ([]*models.Assignment, error)
}
