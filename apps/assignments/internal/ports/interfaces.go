package ports

import (
	"context"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"

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
		tmpl models.AssignmentTemplate,
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

type AssignmentTemplateService interface {
	CreateTemplate(
		ctx context.Context,
		creatorID uuid.UUID,
		dto dto.CreateAssignment,
	) (uuid.UUID, error)
	UpdateTemplate(
		ctx context.Context,
		callerID uuid.UUID,
		id uuid.UUID,
		updates map[string]any,
		targets []dto.Target,
	) (*models.AssignmentTemplate, error)
	DeleteTemplate(
		ctx context.Context,
		callerID uuid.UUID,
		id uuid.UUID,
	) error
	GetTemplate(
		ctx context.Context,
		id uuid.UUID,
	) (*models.AssignmentTemplate, []models.AssignmentTarget, error)
	ListTemplates(
		ctx context.Context,
		creatorID uuid.UUID,
		limit int,
		pageToken string,
	) ([]models.AssignmentTemplateLight, string, error)
}

type SubmissionService interface {
	ListStudentAssignments(
		ctx context.Context,
		studentID uuid.UUID,
		groupID uuid.UUID,
		limit int,
		pageToken string,
		statusFilter domain.SubmissionStatus,
	) ([]dto.StudentItem, string, error)
	StartAssignment(
		ctx context.Context,
		studentID uuid.UUID,
		templateID uuid.UUID,
	) (uuid.UUID, time.Time, error)
	SaveDraft(
		ctx context.Context,
		studentID uuid.UUID,
		dto dto.SaveVersion,
	) (uuid.UUID, error)
	SubmitAssignment(
		ctx context.Context,
		studentID uuid.UUID,
		dto dto.SaveVersion,
	) (uuid.UUID, domain.SubmissionStatus, error)
}

type FeedbackService interface {
	ListSubmissions(
		ctx context.Context,
		templateID uuid.UUID,
		limit int,
		pageToken string,
		filter domain.SubmissionStatus,
	) ([]models.Submission, string, error)
	GetSubmissionDetails(
		ctx context.Context,
		submissionID uuid.UUID,
	) (*dto.FullSubmission, error)
	ProvideFeedback(
		ctx context.Context,
		graderID uuid.UUID,
		dto dto.Feedback,
	) error
}
