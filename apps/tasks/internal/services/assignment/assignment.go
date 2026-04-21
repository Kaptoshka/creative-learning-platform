package assignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"tasks/internal/domain/models"
	"tasks/internal/storage"
)

type AssignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
}

type AssignmentSaver interface {
	SaveAssignment(
		ctx context.Context,
		assignment *models.Assignment,
	) (int64, error)
	UpdateAssignment() error
}

type AssignmentProvider interface {
	AssignmentByID(ctx context.Context, id int64) (*models.Assignment, error)
	AssignmentsByTeacherID(ctx context.Context, teacherID int64) ([]*models.Assignment, error)
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

func (s *AssignmentService) CreateAssignment(
	ctx context.Context,
	assignment *models.Assignment,
) (int64, error) {
	const op = "services.assignment.CreateAssignment"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("creating assignment")

	assignmentID, err := s.assignmentSaver.SaveAssignment(
		ctx,
		assignment,
	)
	if err != nil {
		if errors.Is(err, storage.ErrAssignmentAlreadyExists) {
			log.Warn("assignment already exists")

			return 0, fmt.Errorf("%s: %w", op, storage.ErrAssignmentAlreadyExists)
		}

		log.Error("failed to save assignment")

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("assignment created")

	return assignmentID, nil
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
