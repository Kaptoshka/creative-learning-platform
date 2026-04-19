package assignment

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/auth"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/service"

	"github.com/google/uuid"
)

type assignmentService struct {
	log                *slog.Logger
	assignmentSaver    AssignmentSaver
	assignmentProvider AssignmentProvider
	submissionSaver    SubmissionSaver
	submissionProvider SubmissionProvider
	feedbackSaver      FeedbackSaver
	feedbackProvider   FeedbackProvider
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
	) ([]*models.AssignmentTemplate, error)
}

type SubmissionSaver interface {
	CreateSubmission(
		ctx context.Context,
		submission *models.Submission,
	) error
	AddVersion(
		ctx context.Context,
		version *models.SubmissionVersion,
		isAutosave bool,
	) error
}

type SubmissionProvider interface {
	ByID(
		ctx context.Context,
		submissionID uuid.UUID,
	) (
		*models.Submission,
		error,
	)
	VersionsBySubmissionID(
		ctx context.Context,
		submissionID uuid.UUID,
	) ([]*models.SubmissionVersion, error)
	ListByStudentID(
		ctx context.Context,
		studentID uuid.UUID,
		limit int,
		offset int,
		statusFilter domain.SubmissionStatus,
	) ([]*dto.StudentItem, error)
	ListByAssignmentID(
		ctx context.Context,
		templateID uuid.UUID,
		limit int,
		offset int,
		filter domain.SubmissionStatus,
	) (
		[]*models.Submission,
		error,
	)
}

type FeedbackSaver interface {
	CreateFeedback(
		ctx context.Context,
		feedback *models.Feedback,
		newStatus *domain.SubmissionStatus,
	) error
}

type FeedbackProvider interface {
	BySubmissionID(
		ctx context.Context,
		submissionID uuid.UUID,
	) ([]*models.Feedback, error)
}

func New(
	log *slog.Logger,
	assignmentProvider AssignmentProvider,
	assignmentSaver AssignmentSaver,
	submissionProvider SubmissionProvider,
	submissionSaver SubmissionSaver,
	feedbackProvider FeedbackProvider,
	feedbackSaver FeedbackSaver,
) *assignmentService {
	return &assignmentService{
		log:                log,
		assignmentProvider: assignmentProvider,
		assignmentSaver:    assignmentSaver,
		submissionProvider: submissionProvider,
		submissionSaver:    submissionSaver,
		feedbackProvider:   feedbackProvider,
	}
}

func (s *assignmentService) Create(
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

func (s *assignmentService) Update(
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
		return fmt.Errorf("user is not allowed to delete assignment: %w", domain.ErrForbidden)
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

func (s *assignmentService) List(
	ctx context.Context,
	callerID uuid.UUID,
	limit int,
	pageToken string,
) ([]models.AssignmentTemplateLight, string, error) {
	const op = "services.assignment.List"
	log := s.log.With(slog.String("op", op))

	creatorID := callerID

	offset, err := service.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"caller_id", callerID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = service.NormalizeLimit(limit)

	templates, err := s.assignmentProvider.List(ctx, creatorID, limit, offset)
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

	nextToken := service.EncodePageToken(offset, len(templates), limit)

	log.Info(
		"assignment templates listed successfully",
		"creator_id", creatorID.String(),
		"returned", len(templates),
		"offset", offset,
	)

	return light, nextToken, nil
}

func (s *assignmentService) StartAssignment(
	ctx context.Context,
	studentID uuid.UUID,
	templateID uuid.UUID,
) (uuid.UUID, time.Time, error) {
	const op = "services.assignment.StartAssignment"
	log := s.log.With(slog.String("op", op))

	_, _, err := s.assignmentProvider.ByID(ctx, templateID)
	if err != nil {
		log.Error(
			"failed to get assignment template",
			"template_id", templateID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, time.Time{}, fmt.Errorf("%w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Error(
			"failed to generate submission id",
			"template_id", templateID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, time.Time{}, fmt.Errorf("%w", err)
	}

	now := time.Now().UTC()

	sub := &models.Submission{
		ID:         id,
		TemplateID: templateID,
		StudentID:  studentID,
		Status:     domain.StatusInProgress,
		StartedAt:  now,
	}

	if err := s.submissionSaver.CreateSubmission(ctx, sub); err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Warn(
				"student already started assignment",
				"template_id", templateID.String(),
				"student_id", studentID.String(),
			)
			return uuid.Nil, time.Time{}, fmt.Errorf("%w", domain.ErrAlreadyExists)
		}
		log.Error(
			"failed to create submission",
			"template_id", templateID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, time.Time{}, fmt.Errorf("%w", err)
	}

	log.Info(
		"assignment started successfully",
		"submission_id", sub.ID.String(),
		"template_id", templateID.String(),
		"student_id", studentID.String(),
		"started_at", now,
	)

	return sub.ID, now, nil
}

func (s *assignmentService) SaveDraft(
	ctx context.Context,
	studentID uuid.UUID,
	dto dto.SaveVersion,
) (uuid.UUID, error) {
	const op = "services.assignment.SaveDraft"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionProvider.ByID(ctx, dto.SubmissionID)
	if err != nil {
		log.Error(
			"failed to get submission",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, fmt.Errorf("%w", err)
	}

	if sub.StudentID != studentID {
		log.Warn(
			"student does not own submission",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
			"owner_id", sub.StudentID.String(),
		)
		return uuid.Nil, fmt.Errorf("%w", domain.ErrForbidden)
	}

	if sub.Status != domain.StatusInProgress {
		log.Warn(
			"submission is not in progress",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
			"status", sub.Status,
		)
		return uuid.Nil, fmt.Errorf("%w", domain.ErrSubmissionClosed)
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Error(
			"failed to generate version id",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, fmt.Errorf("%w", err)
	}

	version := &models.SubmissionVersion{
		ID:               id,
		SubmissionID:     dto.SubmissionID,
		Payload:          dto.Payload,
		TimeSpentSeconds: dto.TimeSpent,
		IsAutosave:       true,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.submissionSaver.AddVersion(ctx, version, false); err != nil {
		log.Error(
			"failed to save draft version",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, fmt.Errorf("%w", err)
	}

	log.Info(
		"draft saved successfully",
		"version_id", version.ID.String(),
		"submission_id", dto.SubmissionID.String(),
		"student_id", studentID.String(),
		"time_spent_seconds", version.TimeSpentSeconds,
	)

	return version.ID, nil
}

func (s *assignmentService) SubmitAssignment(
	ctx context.Context,
	studentID uuid.UUID,
	dto dto.SaveVersion,
) (uuid.UUID, domain.SubmissionStatus, error) {
	const op = "services.assignment.SubmitAssignment"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionProvider.ByID(ctx, dto.SubmissionID)
	if err != nil {
		log.Error(
			"failed to get submission",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, "", fmt.Errorf("%w", err)
	}

	if sub.StudentID != studentID {
		log.Warn(
			"student does not own submission",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
			"owner_id", sub.StudentID.String(),
		)
		return uuid.Nil, "", fmt.Errorf("%w", domain.ErrForbidden)
	}

	if sub.Status != domain.StatusInProgress {
		log.Warn(
			"submission is not in progress",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
			"status", sub.Status,
		)
		return uuid.Nil, "", fmt.Errorf("%w", domain.ErrSubmissionClosed)
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Error(
			"failed to generate version ID",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, "", fmt.Errorf("%w", err)
	}

	version := &models.SubmissionVersion{
		ID:               id,
		SubmissionID:     dto.SubmissionID,
		Payload:          dto.Payload,
		TimeSpentSeconds: dto.TimeSpent,
		IsAutosave:       false,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.submissionSaver.AddVersion(ctx, version, true); err != nil {
		log.Error(
			"failed to submit assignment",
			"submission_id", dto.SubmissionID.String(),
			"student_id", studentID.String(),
		)
		return uuid.Nil, "", fmt.Errorf("%w", err)
	}

	log.Info(
		"assignment submitted successfully",
		"version_id", version.ID.String(),
		"submission_id", dto.SubmissionID.String(),
		"student_id", studentID.String(),
		"time_spent_seconds", version.TimeSpentSeconds,
	)

	return version.ID, domain.StatusSubmitted, nil
}

func (s *assignmentService) ListSubmissions(
	ctx context.Context,
	templateID uuid.UUID,
	limit int,
	pageToken string,
	filter domain.SubmissionStatus,
) ([]*models.Submission, string, error) {
	const op = "services.assignment.ListSubmissions"
	log := s.log.With(slog.String("op", op))

	_, _, err := s.assignmentProvider.ByID(ctx, templateID)
	if err != nil {
		log.Error(
			"failed to get assignment template",
			"template_id", templateID.String(),
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	offset, err := service.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"template_id", templateID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = service.NormalizeLimit(limit)

	if err := domain.ValidateSubmissionStatus(filter); err != nil {
		log.Error(
			"invalid status filter",
			"template_id", templateID.String(),
			"filter", filter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	submissions, err := s.submissionProvider.ListByAssignmentID(ctx, templateID, limit, offset, filter)
	if err != nil {
		log.Error(
			"failed to list submissions",
			"template_id", templateID.String(),
			"limit", limit,
			"offset", offset,
			"filter", filter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	nextToken := service.EncodePageToken(offset, len(submissions), limit)

	log.Info(
		"submissions listed successfully",
		"template_id", templateID.String(),
		"returned", len(submissions),
		"offset", offset,
		"filter", filter,
	)

	return submissions, nextToken, nil
}

func (s *assignmentService) ListStudentAssignments(
	ctx context.Context,
	studentID uuid.UUID,
	limit int,
	pageToken string,
	statusFilter domain.SubmissionStatus,
) ([]*dto.StudentItem, string, error) {
	const op = "services.assignment.ListStudentAssignments"
	log := s.log.With(slog.String("op", op))

	offset, err := service.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"student_id", studentID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = service.NormalizeLimit(limit)

	if err := domain.ValidateSubmissionStatus(statusFilter); err != nil {
		log.Error(
			"invalid status filter",
			"student_id", studentID.String(),
			"status_filter", statusFilter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	items, err := s.submissionProvider.ListByStudentID(ctx, studentID, limit, offset, statusFilter)
	if err != nil {
		log.Error(
			"failed to list assignments for student",
			"student_id", studentID.String(),
			"limit", limit,
			"offset", offset,
			"status_filter", statusFilter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	nextToken := service.EncodePageToken(offset, len(items), limit)

	log.Info(
		"student assignments listed successfully",
		"student_id", studentID.String(),
		"returned", len(items),
		"offset", offset,
		"status_filter", statusFilter,
	)

	return items, nextToken, nil
}

func (s *assignmentService) GetSubmissionDetails(
	ctx context.Context,
	submissionID uuid.UUID,
) (*dto.FullSubmission, error) {
	const op = "services.assignment.GetSubmissionDetails"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionProvider.ByID(ctx, submissionID)
	if err != nil {
		log.Error(
			"failed to get submission",
			"submission_id", submissionID.String(),
		)
		return nil, fmt.Errorf("%w", err)
	}

	type templateResult struct {
		tmpl    *models.AssignmentTemplate
		targets []*models.AssignmentTarget
		err     error
	}
	type versionsResult struct {
		versions []*models.SubmissionVersion
		err      error
	}
	type feedbacksResult struct {
		feedbacks []*models.Feedback
		err       error
	}

	tmplCh := make(chan templateResult, 1)
	versionsCh := make(chan versionsResult, 1)
	feedbacksCh := make(chan feedbacksResult, 1)

	go func() {
		tmpl, targets, err := s.assignmentProvider.ByID(ctx, sub.TemplateID)
		tmplCh <- templateResult{tmpl, targets, err}
	}()

	go func() {
		versions, err := s.submissionProvider.VersionsBySubmissionID(ctx, submissionID)
		versionsCh <- versionsResult{versions, err}
	}()

	go func() {
		feedbacks, err := s.feedbackProvider.BySubmissionID(ctx, submissionID)
		feedbacksCh <- feedbacksResult{feedbacks, err}
	}()

	tmplRes := <-tmplCh
	versionsRes := <-versionsCh
	feedbacksRes := <-feedbacksCh

	if tmplRes.err != nil {
		log.Error(
			"failed to get assignment template",
			"submission_id", submissionID.String(),
			"template_id", sub.TemplateID.String(),
		)
		return nil, fmt.Errorf("%w", tmplRes.err)
	}
	if versionsRes.err != nil {
		log.Error(
			"failed to get submission versions",
			"submission_id", submissionID.String(),
		)
		return nil, fmt.Errorf("%w", versionsRes.err)
	}
	if feedbacksRes.err != nil {
		log.Error(
			"failed to get feedbacks",
			"submission_id", submissionID.String(),
		)
		return nil, fmt.Errorf("%w", feedbacksRes.err)
	}

	dto := &dto.FullSubmission{
		Template:   tmplRes.tmpl,
		Targets:    tmplRes.targets,
		Submission: sub,
		Versions:   versionsRes.versions,
		Feedbacks:  feedbacksRes.feedbacks,
	}

	log.Info(
		"submission details fetched successfully",
		"submission_id", submissionID.String(),
		"template_id", sub.TemplateID.String(),
		"versions_count", len(versionsRes.versions),
		"feedbacks_count", len(feedbacksRes.feedbacks),
	)

	return dto, nil
}

func (s *assignmentService) ProvideFeedback(
	ctx context.Context,
	graderID uuid.UUID,
	dto dto.Feedback,
) error {
	const op = "services.assignment.ProvideFeedback"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionProvider.ByID(ctx, dto.SubmissionID)
	if err != nil {
		log.Error(
			"failed to get submission",
			"submission_id", dto.SubmissionID.String(),
			"grader_id", graderID.String(),
		)
		return fmt.Errorf("%w", err)
	}

	if sub.Status == domain.StatusInProgress {
		log.Warn(
			"submission is not yet submitted",
			"submission_id", dto.SubmissionID.String(),
			"grader_id", graderID.String(),
			"status", sub.Status,
		)
		return fmt.Errorf("%w", domain.ErrSubmissionNotSubmitted)
	}

	tmpl, _, err := s.assignmentProvider.ByID(ctx, sub.TemplateID)
	if err != nil {
		log.Error(
			"failed to get assignment template",
			"template_id", sub.TemplateID.String(),
			"grader_id", graderID.String(),
		)
		return fmt.Errorf("%w", err)
	}

	if tmpl.CreatorID != graderID {
		log.Warn(
			"grader does not own assignment template",
			"template_id", sub.TemplateID.String(),
			"grader_id", graderID.String(),
			"owner_id", tmpl.CreatorID.String(),
		)
		return fmt.Errorf("%w", domain.ErrForbidden)
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Error("failed to generate UUIDv7", "error", err)
		return fmt.Errorf("%w", err)
	}

	feedback := &models.Feedback{
		ID:          id,
		VersionID:   dto.VersionID,
		GraderID:    graderID,
		TextContent: &dto.TextContent,
		Payload:     dto.Payload,
		IsPublished: dto.IsPublished,
		CreatedAt:   time.Now().UTC(),
	}

	var newStatus *domain.SubmissionStatus
	if dto.IsPublished && sub.Status == domain.StatusSubmitted {
		status := domain.StatusGraded
		newStatus = &status
	}

	if err := s.feedbackSaver.CreateFeedback(ctx, feedback, newStatus); err != nil {
		log.Error(
			"failed to create feedback",
			"submission_id", dto.SubmissionID.String(),
			"version_id", dto.VersionID.String(),
			"grader_id", graderID.String(),
		)
		return fmt.Errorf("%w", err)
	}

	log.Info(
		"feedback provided successfully",
		"feedback_id", feedback.ID.String(),
		"submission_id", dto.SubmissionID.String(),
		"version_id", dto.VersionID.String(),
		"grader_id", graderID.String(),
		"is_published", dto.IsPublished,
		"new_status", newStatus,
	)

	return nil
}
