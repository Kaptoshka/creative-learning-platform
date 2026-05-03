package grpc

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	assignmentsv1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/assignments/v1"

	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type AssignmentsAdapter struct {
	client assignmentsv1.AssignmentServiceClient
}

func NewAssignmentsAdapter(client assignmentsv1.AssignmentServiceClient) *AssignmentsAdapter {
	return &AssignmentsAdapter{client: client}
}

// --- Teacher: template management ---

func (a *AssignmentsAdapter) CreateAssignment(
	ctx context.Context,
	req domain.CreateAssignmentRequest,
) (domain.CreateAssignmentResponse, error) {
	res, err := a.client.CreateAssignment(ctx, &assignmentsv1.CreateAssignmentRequest{
		Title:        req.Title,
		Description:  req.Description,
		WidgetId:     req.WidgetID,
		WidgetConfig: mapToProto(req.WidgetConfig),
		DueDate:      timeToProto(req.DueDate),
		Targets:      targetsToProto(req.Targets),
	})
	if err != nil {
		return domain.CreateAssignmentResponse{}, mapGRPCError(err)
	}

	return domain.CreateAssignmentResponse{ID: res.Id}, nil
}

func (a *AssignmentsAdapter) UpdateAssignment(
	ctx context.Context,
	req domain.UpdateAssignmentRequest,
) (domain.UpdateAssignmentResponse, error) {
	res, err := a.client.UpdateAssignment(ctx, &assignmentsv1.UpdateAssignmentRequest{
		AssignmentId: req.AssignmentID,
		Template: &assignmentsv1.AssignmentTemplate{
			Id:           req.Template.ID,
			Title:        req.Template.Title,
			Description:  req.Template.Description,
			WidgetId:     req.Template.WidgetID,
			WidgetConfig: mapToProto(req.Template.WidgetConfig),
			DueDate:      timeToProto(req.Template.DueDate),
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: req.UpdateMask},
		Targets:    targetsToProto(req.Targets),
	})
	if err != nil {
		return domain.UpdateAssignmentResponse{}, mapGRPCError(err)
	}

	return domain.UpdateAssignmentResponse{
		Template: protoToTemplate(res.Template),
	}, nil
}

func (a *AssignmentsAdapter) DeleteAssignment(
	ctx context.Context,
	req domain.DeleteAssignmentRequest,
) error {
	_, err := a.client.DeleteAssignment(ctx, &assignmentsv1.DeleteAssignmentRequest{
		Id: req.ID,
	})
	return mapGRPCError(err)
}

func (a *AssignmentsAdapter) GetAssignment(
	ctx context.Context,
	req domain.GetAssignmentRequest,
) (domain.GetAssignmentResponse, error) {
	res, err := a.client.GetAssignment(ctx, &assignmentsv1.GetAssignmentRequest{
		Id: req.ID,
	})
	if err != nil {
		return domain.GetAssignmentResponse{}, mapGRPCError(err)
	}

	return domain.GetAssignmentResponse{
		Template: protoToTemplate(res.Template),
		Targets:  protoToTargets(res.Targets),
	}, nil
}

// --- Teacher: submissions & feedback ---

func (a *AssignmentsAdapter) ListAssignments(
	ctx context.Context,
	req domain.ListAssignmentsRequest,
) (domain.ListAssignmentsResponse, error) {
	res, err := a.client.ListAssignments(ctx, &assignmentsv1.ListAssignmentsRequest{
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
		CreatorId: req.CreatorID,
	})
	if err != nil {
		return domain.ListAssignmentsResponse{}, mapGRPCError(err)
	}

	items := make([]domain.AssignmentTemplateLight, len(res.Items))
	for i, item := range res.Items {
		items[i] = protoToTemplateLight(item)
	}

	return domain.ListAssignmentsResponse{
		Items:         items,
		NextPageToken: res.NextPageToken,
	}, nil
}

func (a *AssignmentsAdapter) ListAssignmentSubmissions(
	ctx context.Context,
	req domain.ListAssignmentSubmissionsRequest,
) (domain.ListAssignmentSubmissionsResponse, error) {
	res, err := a.client.ListAssignmentSubmissions(ctx, &assignmentsv1.ListAssignmentSubmissionsRequest{
		TemplateId:   req.TemplateID,
		PageSize:     req.PageSize,
		PageToken:    req.PageToken,
		StatusFilter: statusToProto(req.StatusFilter),
	})
	if err != nil {
		return domain.ListAssignmentSubmissionsResponse{}, mapGRPCError(err)
	}

	return domain.ListAssignmentSubmissionsResponse{
		Items:         protoToSubmissions(res.Items),
		NextPageToken: res.NextPageToken,
	}, nil
}

func (a *AssignmentsAdapter) GetStudentSubmission(
	ctx context.Context, req domain.GetStudentSubmissionRequest,
) (domain.GetStudentSubmissionResponse, error) {
	res, err := a.client.GetStudentSubmission(ctx, &assignmentsv1.GetStudentSubmissionRequest{
		SubmissionId: req.SubmissionID,
	})
	if err != nil {
		return domain.GetStudentSubmissionResponse{}, mapGRPCError(err)
	}

	return domain.GetStudentSubmissionResponse{
		Template:   protoToTemplate(res.Template),
		Submission: protoToSubmission(res.Submission),
		History:    protoToVersions(res.History),
		Feedback:   protoToFeedbacks(res.Feedback),
	}, nil
}

func (a *AssignmentsAdapter) ProvideFeedback(
	ctx context.Context,
	req domain.ProvideFeedbackRequest,
) error {
	_, err := a.client.ProvideFeedback(ctx, &assignmentsv1.ProvideFeedbackRequest{
		SubmissionId: req.SubmissionID,
		VersionId:    req.VersionID,
		TextContent:  req.TextContent,
		Payload:      mapToProto(req.Payload),
		IsPublished:  req.IsPublished,
	})
	return mapGRPCError(err)
}

// --- Student: assignment workflow ---

func (a *AssignmentsAdapter) ListMyAssignments(
	ctx context.Context,
	req domain.ListMyAssignmentsRequest,
) (domain.ListMyAssignmentsResponse, error) {
	res, err := a.client.ListMyAssignments(ctx, &assignmentsv1.ListMyAssignmentsRequest{
		PageSize:     req.PageSize,
		PageToken:    req.PageToken,
		StatusFilter: statusToProto(req.StatusFilter),
	})
	if err != nil {
		return domain.ListMyAssignmentsResponse{}, mapGRPCError(err)
	}

	items := make([]domain.ListMyAssignmentsItem, len(res.Items))
	for i, item := range res.Items {
		items[i] = domain.ListMyAssignmentsItem{
			Template:    protoToTemplateLight(item.Template),
			Status:      domain.SubmissionStatus(item.Status),
			HasFeedback: item.HasFeedback,
		}
	}

	return domain.ListMyAssignmentsResponse{
		Items:         items,
		NextPageToken: res.NextPageToken,
	}, nil
}

func (a *AssignmentsAdapter) StartAssignment(
	ctx context.Context,
	req domain.StartAssignmentRequest,
) (domain.StartAssignmentResponse, error) {
	res, err := a.client.StartAssignment(ctx, &assignmentsv1.StartAssignmentRequest{
		TemplateId: req.TemplateID,
	})
	if err != nil {
		return domain.StartAssignmentResponse{}, mapGRPCError(err)
	}

	return domain.StartAssignmentResponse{
		SubmissionID: res.SubmissionId,
		StartedAt:    protoToTimeVal(res.StartedAt),
	}, nil
}

func (a *AssignmentsAdapter) SaveAssignmentDraft(
	ctx context.Context,
	req domain.SaveAssignmentDraftRequest,
) (domain.SaveAssignmentDraftResponse, error) {
	res, err := a.client.SaveAssignmentDraft(ctx, &assignmentsv1.SaveAssignmentDraftRequest{
		SubmissionId: req.SubmissionID,
		Payload:      mapToProto(req.Payload),
		TimeSpent:    durationToProto(req.TimeSpent),
	})
	if err != nil {
		return domain.SaveAssignmentDraftResponse{}, mapGRPCError(err)
	}

	return domain.SaveAssignmentDraftResponse{
		VersionID: res.VersionId,
		SavedAt:   protoToTimeVal(res.SavedAt),
	}, nil
}

func (a *AssignmentsAdapter) SubmitAssignment(
	ctx context.Context,
	req domain.SubmitAssignmentRequest,
) (domain.SubmitAssignmentResponse, error) {
	res, err := a.client.SubmitAssignment(ctx, &assignmentsv1.SubmitAssignmentRequest{
		SubmissionId: req.SubmissionID,
		Payload:      mapToProto(req.Payload),
		TimeSpent:    durationToProto(req.TimeSpent),
	})
	if err != nil {
		return domain.SubmitAssignmentResponse{}, mapGRPCError(err)
	}

	return domain.SubmitAssignmentResponse{
		VersionID: res.VersionId,
		Status:    domain.SubmissionStatus(res.Status),
	}, nil
}
