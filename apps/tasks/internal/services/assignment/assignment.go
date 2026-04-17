package assignment

import (
	"log/slog"
	"time"

	"tasks/internal/domain/models"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
)

type AssignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
}

type AssignmentSaver interface {
	CreateAssignment(
		title string,
		description string,
		widgetID string,
		widgetConfig models.JSONB,
		due_date time.Time,
		target *tasksv1.AssignmentTarget,
	) string
	UpdateAssignment(
		assignmentID string,
		template models.AssignmentTemplate,
		updateMask *tasksv1.,
		target,
	)
	DeleteAssignment(

	)
}

type AssignmentProvider interface {
	Assignment(

	)
	ListAssignments(

	)
	ListMyAssignments(

	)
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
