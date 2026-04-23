package driven

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/models"

	"github.com/google/uuid"
)

type AssignmentRepo interface {
	CreateAssignment(
		ctx context.Context,
		tmpl models.AssignmentTemplate,
		targets []models.AssignmentTarget,
	) error
	UpdateAssignment(
		ctx context.Context,
		id uuid.UUID,
		update map[string]any,
		newTargets []models.AssignmentTarget,
	) (*models.AssignmentTemplate, error)
	DeleteAssignment(
		ctx context.Context,
		id uuid.UUID,
	) error
	GetAssignmentByID(
		ctx context.Context,
		id uuid.UUID,
	) (*models.AssignmentTemplate, []models.AssignmentTarget, error)
	ListAssignmentsByCreator(
		ctx context.Context,
		creatorID uuid.UUID,
		limit int,
		offset int,
	) ([]models.AssignmentTemplate, error)
}

type SubmissionRepo interface {
	CreateSubmission(
		ctx context.Context,
		sub models.Submission,
	) error
	GetSubmissionByID(
		ctx context.Context,
		id uuid.UUID,
	) (*models.Submission, error)
	AddSubmissionVersion(
		ctx context.Context,
		version models.SubmissionVersion,
		updateParentStatus bool,
	) error
	GetSubmissionVersions(
		ctx context.Context,
		submissionID uuid.UUID,
	) ([]models.SubmissionVersion, error)
	ListAssignmentsForStudent(
		ctx context.Context,
		studentID uuid.UUID,
		groupID uuid.UUID,
		limit int,
		offset int,
		statusFilter domain.SubmissionStatus,
	) ([]dto.StudentItem, error)
	ListSubmissionsByTemplate(
		ctx context.Context,
		templateID uuid.UUID,
		limit int,
		offset int,
		filter domain.SubmissionStatus,
	) ([]models.Submission, error)
}

type FeedbackRepo interface {
	CreateFeedback(
		ctx context.Context,
		feedback models.Feedback,
		changeSubmissionStatus *domain.SubmissionStatus,
	) error
	GetFeedbacksBySubmission(
		ctx context.Context,
		submissionID uuid.UUID,
	) ([]models.Feedback, error)
}
