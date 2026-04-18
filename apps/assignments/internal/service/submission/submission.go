package submission

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/service"

	"github.com/google/uuid"
)

type SubmissionService struct {
	log                *slog.Logger
	submissionSaver    SubmissionSaver
	submissionProvider SubmissionProvider
}

type SubmissionSaver interface{}

type SubmissionProvider interface {
	ListByStudentID(
		ctx context.Context,
		studentID uuid.UUID,
		limit int,
		offset int,
		statusFilter domain.SubmissionStatus,
	) ([]*dto.StudentItem, int, error)
}

func New(
	log *slog.Logger,
	submissionSaver SubmissionSaver,
	submissionProvider SubmissionProvider,
) *SubmissionService {
	return &SubmissionService{
		log:                log,
		submissionSaver:    submissionSaver,
		submissionProvider: submissionProvider,
	}
}

func (s *SubmissionService) ListStudentAssignments(
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

	items, total, err := s.submissionProvider.ListByStudentID(ctx, studentID, limit, offset, statusFilter)
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

	nextToken := service.EncodePageToken(offset+len(items), total)

	log.Info(
		"student assignments listed successfully",
		"student_id", studentID.String(),
		"returned", len(items),
		"total", total,
		"offset", offset,
		"status_filter", statusFilter,
	)

	return items, nextToken, nil
}
