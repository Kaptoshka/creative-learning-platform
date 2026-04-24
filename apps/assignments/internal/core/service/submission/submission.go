package submission

import (
	"context"
	"errors"
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

type submissionService struct {
	log            *slog.Logger
	assignmentRepo driven.AssignmentRepo
	submissionRepo driven.SubmissionRepo
	feedbackReoi   driven.FeedbackRepo
}

func New(
	log *slog.Logger,
	assignmentRepo driven.AssignmentRepo,
	submissionRepo driven.SubmissionRepo,
	feedbackRepo driven.FeedbackRepo,
) *submissionService {
	return &submissionService{
		log:            log,
		assignmentRepo: assignmentRepo,
		submissionRepo: submissionRepo,
		feedbackReoi:   feedbackRepo,
	}
}

func (s *submissionService) ListStudentAssignments(
	ctx context.Context,
	studentID uuid.UUID,
	limit int,
	pageToken string,
	statusFilter domain.SubmissionStatus,
) ([]dto.StudentItem, string, error) {
	const op = "services.assignment.ListStudentAssignments"
	log := s.log.With(slog.String("op", op))

	groupIDStr := auth.GetGroupID(ctx)
	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidGroupID)
	}

	offset, err := shared.DecodePageToken(pageToken)
	if err != nil {
		log.Error(
			"failed to decode page token",
			"student_id", studentID.String(),
			"group_id", groupID.String(),
			"page_token", pageToken,
		)
		return nil, "", fmt.Errorf("%w", domain.ErrInvalidPageToken)
	}

	limit = shared.NormalizeLimit(limit)

	if err := domain.ValidateSubmissionStatus(statusFilter); err != nil {
		log.Error(
			"invalid status filter",
			"student_id", studentID.String(),
			"group_id", groupID.String(),
			"status_filter", statusFilter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	items, err := s.submissionRepo.ListAssignmentsForStudent(
		ctx, studentID, groupID, limit, offset, statusFilter,
	)
	if err != nil {
		log.Error(
			"failed to list assignments for student",
			"student_id", studentID.String(),
			"group_id", groupID.String(),
			"limit", limit,
			"offset", offset,
			"status_filter", statusFilter,
		)
		return nil, "", fmt.Errorf("%w", err)
	}

	nextToken := shared.EncodePageToken(offset, len(items), limit)

	log.Info(
		"student assignments listed successfully",
		"student_id", studentID.String(),
		"group_id", groupID.String(),
		"returned", len(items),
		"offset", offset,
		"status_filter", statusFilter,
	)

	return items, nextToken, nil
}

func (s *submissionService) StartAssignment(
	ctx context.Context,
	studentID uuid.UUID,
	templateID uuid.UUID,
) (uuid.UUID, time.Time, error) {
	const op = "services.assignment.StartAssignment"
	log := s.log.With(slog.String("op", op))

	_, _, err := s.assignmentRepo.GetAssignmentByID(ctx, templateID)
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

	sub := models.Submission{
		ID:         id,
		TemplateID: templateID,
		StudentID:  studentID,
		Status:     domain.StatusInProgress,
		StartedAt:  now,
	}

	if err := s.submissionRepo.CreateSubmission(ctx, sub); err != nil {
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

func (s *submissionService) SaveDraft(
	ctx context.Context,
	studentID uuid.UUID,
	dto dto.SaveVersion,
) (uuid.UUID, error) {
	const op = "services.assignment.SaveDraft"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionRepo.GetSubmissionByID(ctx, dto.SubmissionID)
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

	version := models.SubmissionVersion{
		ID:               id,
		SubmissionID:     dto.SubmissionID,
		Payload:          dto.Payload,
		TimeSpentSeconds: dto.TimeSpent,
		IsAutosave:       true,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.submissionRepo.AddSubmissionVersion(
		ctx, version, false,
	); err != nil {
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

func (s *submissionService) SubmitAssignment(
	ctx context.Context,
	studentID uuid.UUID,
	dto dto.SaveVersion,
) (uuid.UUID, domain.SubmissionStatus, error) {
	const op = "services.assignment.SubmitAssignment"
	log := s.log.With(slog.String("op", op))

	sub, err := s.submissionRepo.GetSubmissionByID(ctx, dto.SubmissionID)
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

	version := models.SubmissionVersion{
		ID:               id,
		SubmissionID:     dto.SubmissionID,
		Payload:          dto.Payload,
		TimeSpentSeconds: dto.TimeSpent,
		IsAutosave:       false,
		CreatedAt:        time.Now().UTC(),
	}

	if err := s.submissionRepo.AddSubmissionVersion(
		ctx, version, true,
	); err != nil {
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
