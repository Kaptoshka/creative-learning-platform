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
	SubmissionByID(
		ctx context.Context,
		submissionID int64,
	) (models.Submission, error)
	SubmissionsByAssignmentID(
		ctx context.Context,
		assignmentID int64,
	) ([]models.Submission, error)
}
