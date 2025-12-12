package assignment

import (
	"context"
	"fmt"
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
) (models.Assignment, error) {
	const op = "services.assignment.Assignment"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("fetching assignment")

	assignment, err := s.assignmentProvider.AssignmentByID(ctx, assignmentID)
	if err != nil {
		return models.Assignment{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("assignment fetched")

	return *assignment, nil
}
