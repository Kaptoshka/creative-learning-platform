package feedback

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

	"github.com/google/uuid"
)

type feedbackService struct {
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
) *feedbackService {
	return &feedbackService{
		log:            log,
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
		feedbackRepo:   feedbackRepo,
	}
}

func (s *feedbackService) ListSubmissions(
	ctx context.Context,
	templateID uuid.UUID,
	limit int,
	pageToken string,
	filter domain.SubmissionStatus,
) ([]models.Submission, string, error) {
	const op = "services.assignment.ListSubmissions"
	log := s.log.With(slog.String("op", op))

	_, _, err := s.assignmentRepo.GetAssignmentByID(ctx, templateID)
	if err != nil {
		log.Error(
			"failed to get assignment template",
			"template_id", templateID.String(),
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	offset, err := shared.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"template_id", templateID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = shared.NormalizeLimit(limit)

	if err := domain.ValidateSubmissionStatus(filter); err != nil {
		log.Error(
			"invalid status filter",
			"template_id", templateID.String(),
			"filter", filter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	submissions, err := s.submissionRepo.ListSubmissionsByTemplate(
		ctx, templateID, limit, offset, filter,
	)
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

	nextToken := shared.EncodePageToken(offset, len(submissions), limit)

	log.Info(
		"submissions listed successfully",
		"template_id", templateID.String(),
		"returned", len(submissions),
		"offset", offset,
		"filter", filter,
	)

	return submissions, nextToken, nil
}

func (s *feedbackService) GetSubmissionDetails(
	ctx context.Context,
	submissionID uuid.UUID,
) (*dto.FullSubmission, error) {
	const op = "services.assignment.GetSubmissionDetails"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionRepo.GetSubmissionByID(ctx, submissionID)
	if err != nil {
		log.Error(
			"failed to get submission",
			"submission_id", submissionID.String(),
		)
		return nil, fmt.Errorf("%w", err)
	}

	type templateResult struct {
		tmpl    *models.AssignmentTemplate
		targets []models.AssignmentTarget
		err     error
	}
	type versionsResult struct {
		versions []models.SubmissionVersion
		err      error
	}
	type feedbacksResult struct {
		feedbacks []models.Feedback
		err       error
	}

	tmplCh := make(chan templateResult, 1)
	versionsCh := make(chan versionsResult, 1)
	feedbacksCh := make(chan feedbacksResult, 1)

	go func() {
		tmpl, targets, err := s.assignmentRepo.GetAssignmentByID(
			ctx, sub.TemplateID,
		)
		tmplCh <- templateResult{tmpl, targets, err}
	}()

	go func() {
		versions, err := s.submissionRepo.GetSubmissionVersions(
			ctx, submissionID,
		)
		versionsCh <- versionsResult{versions, err}
	}()

	go func() {
		feedbacks, err := s.feedbackRepo.GetFeedbacksBySubmission(
			ctx, submissionID,
		)
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
		Assignment: tmplRes.tmpl,
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

func (s *feedbackService) ProvideFeedback(
	ctx context.Context,
	graderID uuid.UUID,
	dto dto.Feedback,
) error {
	const op = "services.assignment.ProvideFeedback"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionRepo.GetSubmissionByID(
		ctx, dto.SubmissionID,
	)
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

	tmpl, _, err := s.assignmentRepo.GetAssignmentByID(ctx, sub.TemplateID)
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

	feedback := models.Feedback{
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

	if err := s.feedbackRepo.CreateFeedback(
		ctx, feedback, newStatus,
	); err != nil {
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
