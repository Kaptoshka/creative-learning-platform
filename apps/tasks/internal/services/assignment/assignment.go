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

func New(
	log *slog.Logger,
	assignmentProvider AssignmentProvider,
	assignmentSaver AssignmentSaver,
) *AssignmentService {
	return &AssignmentService{
		log:                log,
		assignmentProvider: assignmentProvider,
		assignmentSaver:    assignmentSaver,
	}
}

func (s *AssignmentService) Assignment(
	ctx context.Context,
	assignmentID int64,
) (models.Assignment, models.Submission, error) {

}
