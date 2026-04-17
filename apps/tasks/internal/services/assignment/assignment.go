package assignment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tasks/internal/domain/models"
	"tasks/internal/services"

	"github.com/google/uuid"
)

const (
	cutoffDuration = 5 * time.Minute
)

type AssignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
}

type AssignmentSaver interface {
	CreateAssignment(
		ctx context.Context,
		template models.AssignmentTemplate,
		targets []models.AssignmentTarget,
	) error
	UpdateAssignment(
		assignmentID string,
		template models.AssignmentTemplate,
		// TODO: cannot be pb type
		// target,
	)
	DeleteAssignment()
}

type AssignmentProvider interface {
	Assignment()
	ListAssignments()
	ListMyAssignments()
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

func (s *AssignmentService) Create(
	ctx context.Context,
	creatorID uuid.UUID,
	dto models.CreateAssignmentDTO,
) (uuid.UUID, error) {
	const op = "services.assignment.Create"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("validation of the assignment DTO")

	if err := validateCreateTemplateDTO(dto); err != nil {
		return uuid.Nil, fmt.Errorf("validation error: %w", err)
	}

	now := time.Now().UTC()

	tmplID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("uuid generation error: %w", err)
	}

	tmpl := models.AssignmentTemplate{
		ID:           tmplID,
		CreatorID:    creatorID,
		Title:        dto.Title,
		Description:  dto.Description,
		WidgetID:     dto.WidgetID,
		WidgetConfig: dto.WidgetConfig,
		DueDate:      dto.DueDate,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	targets := make([]models.AssignmentTarget, 0, len(dto.Targets))
	var targetID uuid.UUID
	for _, t := range dto.Targets {
		targetID, err = uuid.NewV7()
		if err != nil {
			return uuid.Nil, fmt.Errorf("uuid generation error: %w", err)
		}
		targets = append(targets, models.AssignmentTarget{
			ID:         targetID,
			TemplateID: tmpl.ID,
			GroupID:    t.GroupID,
			StudentID:  t.StudentID,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}

	if err := s.assignmentSaver.CreateAssignment(ctx, tmpl, targets); err != nil {
		return uuid.Nil, fmt.Errorf("storage error: %w", err)
	}

	log.Debug("assignment created", "id", tmpl.ID)

	return tmpl.ID, nil
}

func validateCreateTemplateDTO(dto models.CreateAssignmentDTO) error {
	if dto.Title == "" {
		return services.ErrTitleRequired
	}
	if dto.WidgetID == uuid.Nil {
		return services.ErrWidgetIDRequired
	}

	if dto.DueDate != nil {
		deadline := dto.DueDate.Add(-cutoffDuration)
		if !deadline.After(time.Now().UTC()) {
			return services.ErrDueDateTooClose
		}
	}

	if len(dto.Targets) == 0 {
		return services.ErrTargetsRequired
	}
	for i, t := range dto.Targets {
		if t.GroupID == nil && t.StudentID == nil {
			return fmt.Errorf("%w: target[%d]", services.ErrTargetEmpty, i)
		}
	}

	return nil
}
