package usecase

import (
	"context"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type AssignmentsUseCase struct {
	assignments outgoing.AssignmentService
}

func NewAssignmentsUseCase(
	assignments outgoing.AssignmentService,
) *AssignmentsUseCase {
	return &AssignmentsUseCase{assignments: assignments}
}

// --- Teacher: template management ---

func (uc *AssignmentsUseCase) CreateAssignment(
	ctx context.Context,
	req domain.CreateAssignmentRequest,
) (domain.CreateAssignmentResponse, error) {
	if req.Title == "" {
		return domain.CreateAssignmentResponse{}, fmt.Errorf(
			"%w: title is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.WidgetID == "" {
		return domain.CreateAssignmentResponse{}, fmt.Errorf(
			"%w: widget_id is required",
			domain.ErrInvalidArgument,
		)
	}
	for _, t := range req.Targets {
		if err := validateTarget(t); err != nil {
			return domain.CreateAssignmentResponse{}, err
		}
	}

	return uc.assignments.CreateAssignment(ctx, req)
}

func (uc *AssignmentsUseCase) UpdateAssignment(
	ctx context.Context,
	req domain.UpdateAssignmentRequest,
) (domain.UpdateAssignmentResponse, error) {
	if req.AssignmentID == "" {
		return domain.UpdateAssignmentResponse{}, fmt.Errorf(
			"%w: assignment_id is required",
			domain.ErrInvalidArgument,
		)
	}
	if len(req.UpdateMask) == 0 {
		return domain.UpdateAssignmentResponse{}, fmt.Errorf(
			"%w: update_mask is required",
			domain.ErrInvalidArgument,
		)
	}
	for _, t := range req.Targets {
		if err := validateTarget(t); err != nil {
			return domain.UpdateAssignmentResponse{}, err
		}
	}

	return uc.assignments.UpdateAssignment(ctx, req)
}

func (uc *AssignmentsUseCase) DeleteAssignment(
	ctx context.Context,
	req domain.DeleteAssignmentRequest,
) error {
	if req.ID == "" {
		return fmt.Errorf("%w: id is required", domain.ErrInvalidArgument)
	}

	return uc.assignments.DeleteAssignment(ctx, req)
}

func (uc *AssignmentsUseCase) GetAssignment(
	ctx context.Context,
	req domain.GetAssignmentRequest,
) (domain.GetAssignmentResponse, error) {
	if req.ID == "" {
		return domain.GetAssignmentResponse{}, fmt.Errorf(
			"%w: id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.GetAssignment(ctx, req)
}

// --- Teacher: submissions & feedback ---

func (uc *AssignmentsUseCase) ListAssignments(
	ctx context.Context,
	req domain.ListAssignmentsRequest,
) (domain.ListAssignmentsResponse, error) {
	return uc.assignments.ListAssignments(ctx, req)
}

func (uc *AssignmentsUseCase) ListAssignmentSubmissions(
	ctx context.Context,
	req domain.ListAssignmentSubmissionsRequest,
) (domain.ListAssignmentSubmissionsResponse, error) {
	if req.TemplateID == "" {
		return domain.ListAssignmentSubmissionsResponse{}, fmt.Errorf(
			"%w: template_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.ListAssignmentSubmissions(ctx, req)
}

func (uc *AssignmentsUseCase) GetStudentSubmission(
	ctx context.Context,
	req domain.GetStudentSubmissionRequest,
) (domain.GetStudentSubmissionResponse, error) {
	if req.SubmissionID == "" {
		return domain.GetStudentSubmissionResponse{}, fmt.Errorf(
			"%w: submission_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.GetStudentSubmission(ctx, req)
}

func (uc *AssignmentsUseCase) ProvideFeedback(
	ctx context.Context,
	req domain.ProvideFeedbackRequest,
) error {
	if req.SubmissionID == "" {
		return fmt.Errorf("%w: submission_id is required", domain.ErrInvalidArgument)
	}
	if req.VersionID == "" {
		return fmt.Errorf("%w: version_id is required", domain.ErrInvalidArgument)
	}
	if req.TextContent == "" {
		return fmt.Errorf("%w: text_content is required", domain.ErrInvalidArgument)
	}

	return uc.assignments.ProvideFeedback(ctx, req)
}

// --- Student: assignment workflow ---

func (uc *AssignmentsUseCase) ListMyAssignments(
	ctx context.Context,
	req domain.ListMyAssignmentsRequest,
) (domain.ListMyAssignmentsResponse, error) {
	return uc.assignments.ListMyAssignments(ctx, req)
}

func (uc *AssignmentsUseCase) StartAssignment(
	ctx context.Context,
	req domain.StartAssignmentRequest,
) (domain.StartAssignmentResponse, error) {
	if req.TemplateID == "" {
		return domain.StartAssignmentResponse{}, fmt.Errorf(
			"%w: template_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.StartAssignment(ctx, req)
}

func (uc *AssignmentsUseCase) SaveAssignmentDraft(
	ctx context.Context,
	req domain.SaveAssignmentDraftRequest,
) (domain.SaveAssignmentDraftResponse, error) {
	if req.SubmissionID == "" {
		return domain.SaveAssignmentDraftResponse{}, fmt.Errorf(
			"%w: submission_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.SaveAssignmentDraft(ctx, req)
}

func (uc *AssignmentsUseCase) SubmitAssignment(
	ctx context.Context,
	req domain.SubmitAssignmentRequest,
) (domain.SubmitAssignmentResponse, error) {
	if req.SubmissionID == "" {
		return domain.SubmitAssignmentResponse{}, fmt.Errorf(
			"%w: submission_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.assignments.SubmitAssignment(ctx, req)
}

// --- Helpers ---

func validateTarget(t domain.AssignmentTarget) error {
	if t.GroupID == "" && t.StudentID == "" {
		return fmt.Errorf(
			"%w: target must have group_id or student_id",
			domain.ErrInvalidArgument,
		)
	}
	if t.GroupID != "" && t.StudentID != "" {
		return fmt.Errorf(
			"%w: target must have either group_id or student_id, not both",
			domain.ErrInvalidArgument,
		)
	}
	return nil
}
