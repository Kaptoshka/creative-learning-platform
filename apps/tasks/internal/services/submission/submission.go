package submission

import (
	"context"
	"encoding/json"
	"log/slog"
	"tasks/internal/domain/models"
)

type SubmissionService struct {
	log                *slog.Logger
	submissionSaver    SubmissionSaver
	submissionProvider SubmissionProvider
}

type SubmissionSaver interface {
	SaveSubmission(
		ctx context.Context,
		assignmentID int64,
		studentID int64,
		content json.RawMessage,
	) (int64, error)
	UpdateSubmission(
		ctx context.Context,
		submissionID int64,
		content json.RawMessage,
	) error
}

type SubmissionProvider interface {
	SubmissionByAssignmentID(
		ctx context.Context,
		assignmentID int64,
	) ([]models.Submission, error)
	SubmissionsByAssignmentIDAndStudentID(
		ctx context.Context,
		assignmentID int64,
		studentID int64,
	) (models.Submission, error)
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

func (s *SubmissionService) SubmissionByAssignmentIDAndStudentID(
	ctx context.Context,
	assignmentID int64,
	studentID int64,
) (models.Submission, error) {
	const op = "services.submission.SubmissionByAssignmentIDAndStudentID"

	log := s.log.With(
		slog.String("op", op),
	)

	log.Debug("fetching submission by assignment and student id's")

	submission, err := s.submissionProvider.SubmissionsByAssignmentIDAndStudentID(
		ctx, assignmentID, studentID,
	)
	if err != nil {
		return models.Submission{}, err
	}

	return submission, nil
}
