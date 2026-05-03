package outgoing

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
)

type AssignmentService interface {
	// Teacher: template management
	CreateAssignment(
		ctx context.Context,
		req domain.CreateAssignmentRequest,
	) (domain.CreateAssignmentResponse, error)
	UpdateAssignment(
		ctx context.Context,
		req domain.UpdateAssignmentRequest,
	) (domain.UpdateAssignmentResponse, error)
	DeleteAssignment(
		ctx context.Context,
		req domain.DeleteAssignmentRequest,
	) error
	GetAssignment(
		ctx context.Context,
		req domain.GetAssignmentRequest,
	) (domain.GetAssignmentResponse, error)

	// Teacher: submissions & feedback
	ListAssignments(
		ctx context.Context,
		req domain.ListAssignmentsRequest,
	) (domain.ListAssignmentsResponse, error)
	ListAssignmentSubmissions(
		ctx context.Context,
		req domain.ListAssignmentSubmissionsRequest,
	) (domain.ListAssignmentSubmissionsResponse, error)
	GetStudentSubmission(
		ctx context.Context,
		req domain.GetStudentSubmissionRequest,
	) (domain.GetStudentSubmissionResponse, error)
	ProvideFeedback(
		ctx context.Context,
		req domain.ProvideFeedbackRequest,
	) error

	// Student: assignment workflow
	ListMyAssignments(
		ctx context.Context,
		req domain.ListMyAssignmentsRequest,
	) (domain.ListMyAssignmentsResponse, error)
	StartAssignment(
		ctx context.Context,
		req domain.StartAssignmentRequest,
	) (domain.StartAssignmentResponse, error)
	SaveAssignmentDraft(
		ctx context.Context,
		req domain.SaveAssignmentDraftRequest,
	) (domain.SaveAssignmentDraftResponse, error)
	SubmitAssignment(
		ctx context.Context,
		req domain.SubmitAssignmentRequest,
	) (domain.SubmitAssignmentResponse, error)
}
