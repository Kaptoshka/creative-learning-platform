package driving

import (
	"context"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/models"

	"github.com/google/uuid"
)

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
