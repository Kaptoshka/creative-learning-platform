package assignment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tasks/internal/auth"
	"tasks/internal/domain/models"
	"tasks/internal/services"

	"github.com/google/uuid"
)

const (
	cutoffDuration = 5 * time.Minute
)

type assignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
}

type AssignmentSaver interface {
	CreateAssignment(
		ctx context.Context,
		template *models.AssignmentTemplate,
		targets []*models.AssignmentTarget,
	) error
	UpdateAssignment(
		ctx context.Context,
		id uuid.UUID,
		updates map[string]any,
		newTargets []*models.AssignmentTarget,
	) (*models.AssignmentTemplate, error)
	DeleteAssignment(
		ctx context.Context,
		id uuid.UUID,
	) error
}

type AssignmentProvider interface {
	ByID(
		ctx context.Context,
		assignmentID uuid.UUID,
	) (
		*models.AssignmentTemplate,
		[]*models.AssignmentTarget,
		error,
	)
	List(
		ctx context.Context,
		creatorID uuid.UUID,
		limit int,
		offset int,
	) (
		[]*models.AssignmentTemplate,
		int,
		error,
	)
}

func New(
	log *slog.Logger,
	assignmentProvider AssignmentProvider,
	assignmentSaver AssignmentSaver,
) *assignmentService {
	return &assignmentService{
		log:                log,
		assignmentProvider: assignmentProvider,
		assignmentSaver:    assignmentSaver,
	}
}

func (s *assignmentService) Create(
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

	tmpl := &models.AssignmentTemplate{
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

	targets := make([]*models.AssignmentTarget, 0, len(dto.Targets))
	var targetID uuid.UUID
	for _, t := range dto.Targets {
		targetID, err = uuid.NewV7()
		if err != nil {
			return uuid.Nil, fmt.Errorf("uuid generation error: %w", err)
		}
		targets = append(targets, &models.AssignmentTarget{
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

func (s *assignmentService) Update(
	ctx context.Context,
	callerID uuid.UUID,
	id uuid.UUID,
	updates map[string]any,
	targets []models.TargetDTO,
) (*models.AssignmentTemplate, error) {
	const op = "services.assignment.UpdateAssignment"

	log := s.log.With(
		slog.String("op", op),
	)

	tmpl, _, err := s.assignmentProvider.ByID(ctx, id)
	if err != nil {
		log.Error(
			"failed to fetch assignment",
			"id",
			id.String(),
		)
		return nil, fmt.Errorf("failed to fetch assignment: %w", err)
	}

	role := auth.GetUserRole(ctx)
	if tmpl.CreatorID != callerID && role != auth.RoleAdmin {
		log.Warn(
			"user is not allowed to update this assignment",
			"id",
			id.String(),
		)
		return nil, fmt.Errorf(
			"user is not allowed to update this assignment: %w",
			services.ErrForbidden,
		)
	}

	if err := validateTemplateUpdates(updates); err != nil {
		log.Error(
			"failed validate updates",
			"id",
			id.String(),
		)
		return nil, fmt.Errorf("failed validate updates: %w", err)
	}

	var domainTargets []*models.AssignmentTarget
	if len(targets) > 0 {
		if err := validateTargets(targets); err != nil {
			log.Error(
				"failed to validate targets",
				"id",
				id.String(),
			)
			return nil, fmt.Errorf("failed to validate targets: %w", err)
		}
		domainTargets = convertTargets(id, targets)
	}

	updates["updated_at"] = time.Now().UTC()

	updated, err := s.assignmentSaver.UpdateAssignment(
		ctx, id, updates, domainTargets,
	)
	if err != nil {
		log.Error(
			"failed to update assignment",
			"id",
			id.String(),
		)
		return nil, fmt.Errorf("failed to update assignment: %w", err)
	}

	return updated, nil
}

func validateTemplateUpdates(updates map[string]any) error {
	if len(updates) == 0 {
		return services.ErrNoUpdates
	}

	if v, ok := updates["title"]; ok {
		title, ok := v.(string)
		if !ok || title == "" {
			return services.ErrTitleRequired
		}
	}

	if v, ok := updates["widget_id"]; ok {
		wid, ok := v.(uuid.UUID)
		if !ok || wid == uuid.Nil {
			return services.ErrWidgetIDRequired
		}
	}

	if v, ok := updates["due_date"]; ok && v != nil {
		dueDate, ok := v.(time.Time)
		if !ok {
			return services.ErrInvalidDueDate
		}
		deadline := dueDate.Add(-cutoffDuration)
		if !deadline.After(time.Now().UTC()) {
			return services.ErrDueDateTooClose
		}
	}

	for _, forbidden := range []string{"id", "creator_id", "created_at"} {
		if _, ok := updates[forbidden]; ok {
			return fmt.Errorf("%w: %s", services.ErrForbiddenField, forbidden)
		}
	}

	return nil
}

func validateTargets(targets []models.TargetDTO) error {
	for i, t := range targets {
		if t.GroupID == nil && t.StudentID == nil {
			return fmt.Errorf("%w: target[%d]", services.ErrTargetEmpty, i)
		}
	}
	return nil
}

func convertTargets(
	templateID uuid.UUID,
	targets []models.TargetDTO,
) []*models.AssignmentTarget {
	now := time.Now().UTC()
	result := make(
		[]*models.AssignmentTarget,
		0,
		len(targets),
	)
	var id uuid.UUID
	var err error
	for _, t := range targets {
		id, err = uuid.NewV7()
		if err != nil {
			return nil
		}
		result = append(result, &models.AssignmentTarget{
			ID:         id,
			TemplateID: templateID,
			GroupID:    t.GroupID,
			StudentID:  t.StudentID,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}
	return result
}

func (s *assignmentService) DeleteAssignment(
	ctx context.Context,
	callerID uuid.UUID,
	id uuid.UUID,
) error {
	const op = "services.assignment.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("attempting to delete assignment")

	tmpl, _, err := s.assignmentProvider.ByID(ctx, id)
	if err != nil {
		log.Error(
			"failed to fetch assignment",
			"id",
			id.String(),
			"error",
			err.Error(),
		)
		return fmt.Errorf("failed to fetch assignment: %w", err)
	}

	role := auth.GetUserRole(ctx)
	if tmpl.CreatorID != callerID && role != auth.RoleAdmin {
		log.Warn(
			"user is not allowed to delete assignment",
			"caller_id",
			callerID.String(),
			"assignment_id",
			id.String(),
		)
		return fmt.Errorf("user is not allowed to delete assignment: %w", services.ErrForbidden)
	}

	if err := s.assignmentSaver.DeleteAssignment(ctx, id); err != nil {
		log.Error(
			"failed to delete assignment",
			"id",
			id.String(),
		)
		return fmt.Errorf("failed to delete assignment: %w", err)
	}

	log.Debug("assignment deleted", "id", id.String())

	return nil
}

func (s *assignmentService) GetTemplate(
	ctx context.Context,
	id uuid.UUID,
) (
	*models.AssignmentTemplate,
	[]*models.AssignmentTarget,
	error,
) {
	const op = "services.assignment.GetTemplate"
	log := s.log.With(slog.String("op", op))

	tmpl, targets, err := s.assignmentProvider.ByID(ctx, id)
	if err != nil {
		log.Error(
			"failed to get assignment template",
			"id", id.String(),
		)
		return nil, nil, fmt.Errorf("%w", err)
	}

	log.Info(
		"assignment template fetched successfully",
		"id", id.String(),
		"targets_count", len(targets),
	)

	return tmpl, targets, nil
}
