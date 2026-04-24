package assignment

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/models"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/shared"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/ports/driven"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/pkg/auth"

	"github.com/google/uuid"
)

type assignmentService struct {
	log            *slog.Logger
	assignmentRepo driven.AssignmentRepo
	submissionRepo driven.SubmissionRepo
	feedbackRepo   driven.FeedbackRepo
}

func New(
	log *slog.Logger,
	assignmentRepo driven.AssignmentRepo,
	submissionRepo driven.SubmissionRepo,
	feedbackRepo driven.FeedbackRepo,
) *assignmentService {
	return &assignmentService{
		log:            log,
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
		feedbackRepo:   feedbackRepo,
	}
}

func (s *assignmentService) CreateTemplate(
	ctx context.Context,
	creatorID uuid.UUID,
	dto dto.CreateAssignment,
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

	if err := s.assignmentRepo.CreateAssignment(ctx, tmpl, targets); err != nil {
		return uuid.Nil, fmt.Errorf("storage error: %w", err)
	}

	log.Debug("assignment created", "id", tmpl.ID)

	return tmpl.ID, nil
}

func validateCreateTemplateDTO(dto dto.CreateAssignment) error {
	if dto.Title == "" {
		return domain.ErrTitleRequired
	}
	if dto.WidgetID == uuid.Nil {
		return domain.ErrWidgetIDRequired
	}

	if dto.DueDate != nil {
		deadline := dto.DueDate.Add(-domain.CutoffDuration)
		if !deadline.After(time.Now().UTC()) {
			return domain.ErrDueDateTooClose
		}
	}

	if len(dto.Targets) == 0 {
		return domain.ErrTargetsRequired
	}
	for i, t := range dto.Targets {
		if t.GroupID == nil && t.StudentID == nil {
			return fmt.Errorf("%w: target[%d]", domain.ErrTargetEmpty, i)
		}
	}

	return nil
}

func (s *assignmentService) UpdateTemplate(
	ctx context.Context,
	callerID uuid.UUID,
	id uuid.UUID,
	updates map[string]any,
	targets []dto.Target,
) (*models.AssignmentTemplate, error) {
	const op = "services.assignment.UpdateAssignment"

	log := s.log.With(
		slog.String("op", op),
	)

	tmpl, _, err := s.assignmentRepo.GetAssignmentByID(ctx, id)
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
			domain.ErrForbidden,
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

	var domainTargets []models.AssignmentTarget
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

	updated, err := s.assignmentRepo.UpdateAssignment(
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
		return domain.ErrNoUpdates
	}

	if v, ok := updates["title"]; ok {
		title, ok := v.(string)
		if !ok || title == "" {
			return domain.ErrTitleRequired
		}
	}

	if v, ok := updates["widget_id"]; ok {
		wid, ok := v.(uuid.UUID)
		if !ok || wid == uuid.Nil {
			return domain.ErrWidgetIDRequired
		}
	}

	if v, ok := updates["due_date"]; ok && v != nil {
		dueDate, ok := v.(time.Time)
		if !ok {
			return domain.ErrInvalidDueDate
		}
		deadline := dueDate.Add(-domain.CutoffDuration)
		if !deadline.After(time.Now().UTC()) {
			return domain.ErrDueDateTooClose
		}
	}

	for _, forbidden := range []string{"id", "creator_id", "created_at"} {
		if _, ok := updates[forbidden]; ok {
			return fmt.Errorf("%w: %s", domain.ErrForbiddenField, forbidden)
		}
	}

	return nil
}

func validateTargets(targets []dto.Target) error {
	for i, t := range targets {
		if t.GroupID == nil && t.StudentID == nil {
			return fmt.Errorf("%w: target[%d]", domain.ErrTargetEmpty, i)
		}
	}
	return nil
}

func convertTargets(
	templateID uuid.UUID,
	targets []dto.Target,
) []models.AssignmentTarget {
	now := time.Now().UTC()
	result := make(
		[]models.AssignmentTarget,
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
		result = append(result, models.AssignmentTarget{
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

func (s *assignmentService) DeleteTemplate(
	ctx context.Context,
	callerID uuid.UUID,
	id uuid.UUID,
) error {
	const op = "services.assignment.Delete"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("attempting to delete assignment")

	tmpl, _, err := s.assignmentRepo.GetAssignmentByID(ctx, id)
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
		return fmt.Errorf(
			"user is not allowed to delete assignment: %w",
			domain.ErrForbidden,
		)
	}

	if err := s.assignmentRepo.DeleteAssignment(ctx, id); err != nil {
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
	[]models.AssignmentTarget,
	error,
) {
	const op = "services.assignment.GetTemplate"
	log := s.log.With(slog.String("op", op))

	tmpl, targets, err := s.assignmentRepo.GetAssignmentByID(ctx, id)
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

func (s *assignmentService) ListTemplates(
	ctx context.Context,
	callerID uuid.UUID,
	limit int,
	pageToken string,
) ([]models.AssignmentTemplateLight, string, error) {
	const op = "services.assignment.List"
	log := s.log.With(slog.String("op", op))

	creatorID := callerID

	offset, err := shared.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"caller_id", callerID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = shared.NormalizeLimit(limit)

	templates, err := s.assignmentRepo.ListTemplatesByCreator(
		ctx, creatorID, limit, offset,
	)
	if err != nil {
		log.Error(
			"failed to list assignment templates",
			"creator_id", creatorID.String(),
			"limit", limit,
			"offset", offset,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	light := make([]models.AssignmentTemplateLight, 0, len(templates))
	for _, t := range templates {
		light = append(light, models.AssignmentTemplateLight{
			ID:      t.ID,
			Title:   t.Title,
			DueDate: t.DueDate,
		})
	}

	nextToken := shared.EncodePageToken(offset, len(templates), limit)

	log.Info(
		"assignment templates listed successfully",
		"creator_id", creatorID.String(),
		"returned", len(templates),
		"offset", offset,
	)

	return light, nextToken, nil
}
