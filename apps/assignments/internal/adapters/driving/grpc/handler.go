package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/ports/driving"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/pkg/auth"

	assignmentsv1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/proto/assignments/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type serverAPI struct {
	assignmentsv1.UnimplementedAssignmentServiceServer
	assignments driving.AssignmentService
	submissions driving.SubmissionService
	feedbacks   driving.FeedbackService
}

func Register(
	gRPC *grpc.Server,
	assignments driving.AssignmentService,
	submissions driving.SubmissionService,
	feedbacks driving.FeedbackService,
) {
	assignmentsv1.RegisterAssignmentServiceServer(gRPC, &serverAPI{
		assignments: assignments,
		submissions: submissions,
		feedbacks:   feedbacks,
	})
}

func (s *serverAPI) CreateAssignment(
	ctx context.Context,
	req *assignmentsv1.CreateAssignmentRequest,
) (*assignmentsv1.CreateAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID in token")
	}

	targets := make([]*dto.Target, 0, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid target")
		}

		targets = append(targets, &target)
	}

	widgetID, err := uuid.Parse(req.WidgetId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "parsing widget ID error")
	}

	widgetConfig, err := structPBToRawMessage(req.GetWidgetConfig())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "parsing widget config error")
	}

	dueTime := req.DueDate.AsTime()

	assignmentDTO := dto.CreateAssignment{
		Title:        req.Title,
		Description:  req.Description,
		WidgetID:     widgetID,
		WidgetConfig: widgetConfig,
		DueDate:      &dueTime,
		Targets:      targets,
	}

	assignmentID, err := s.assignments.CreateTemplate(
		ctx,
		userID,
		assignmentDTO,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "creating assignment error")
	}

	return &assignmentsv1.CreateAssignmentResponse{
		Id: assignmentID.String(),
	}, nil
}

func (s *serverAPI) UpdateAssignment(
	ctx context.Context,
	req *assignmentsv1.UpdateAssignmentRequest,
) (*assignmentsv1.UpdateAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user id",
		)
	}

	userRole := auth.GetUserRole(ctx)
	if userRole != auth.RoleAdmin && userRole != auth.RoleTeacher {
		return nil, status.Error(
			codes.PermissionDenied,
			"user not allowed to update assignment",
		)
	}

	if !req.UpdateMask.IsValid(req.Template) {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid update mask",
		)
	}

	updates := make(map[string]any)

	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "title":
			updates["title"] = req.Template.Title
		case "description":
			updates["description"] = req.Template.Description
		case "widget_id":
			updates["widget_id"] = req.Template.WidgetId
		case "widget_config":
			updates["widget_config"] = req.Template.WidgetConfig
		case "due_date":
			updates["due_date"] = req.Template.DueDate.AsTime().Unix()
		}
	}

	targets := make([]dto.Target, 0, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid target",
			)
		}

		targets = append(targets, target)
	}

	assignmentID, err := uuid.Parse(req.AssignmentId)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"assignment ID cannot be parsed",
		)
	}

	updateModel, err := s.assignments.UpdateTemplate(
		ctx, userID, assignmentID, updates, targets,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"cannot update assignment",
		)
	}

	widgetConfig, err := rawMessageToStructPB(updateModel.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &assignmentsv1.UpdateAssignmentResponse{
		Template: &assignmentsv1.AssignmentTemplate{
			Id:           updateModel.ID.String(),
			CreatorId:    updateModel.CreatorID.String(),
			Title:        updateModel.Title,
			Description:  updateModel.Description,
			WidgetId:     updateModel.WidgetID.String(),
			WidgetConfig: widgetConfig,
			DueDate:      timestamppb.New(*updateModel.DueDate),
			CreatedAt:    timestamppb.New(updateModel.CreatedAt),
			UpdatedAt:    timestamppb.New(updateModel.UpdatedAt),
		},
	}, nil
}

func processTarget(
	t *assignmentsv1.AssignmentTarget,
) (dto.Target, error) {
	switch v := t.GetTarget().(type) {
	case *assignmentsv1.AssignmentTarget_GroupId:
		groupID, err := uuid.Parse(v.GroupId)
		if err != nil {
			return dto.Target{}, errors.New("invalid group ID")
		}

		return dto.Target{
			GroupID: &groupID,
		}, nil
	case *assignmentsv1.AssignmentTarget_StudentId:
		studentID, err := uuid.Parse(v.StudentId)
		if err != nil {
			return dto.Target{}, errors.New("invalid student ID")
		}

		return dto.Target{
			StudentID: &studentID,
		}, nil

	case nil:
		return dto.Target{}, nil

	default:
		return dto.Target{}, errors.New("unknown target type")
	}
}

func (s *serverAPI) DeleteAssignment(
	ctx context.Context,
	req *assignmentsv1.DeleteAssignmentRequest,
) (*assignmentsv1.DeleteAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "user not authenticated")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed parse uuid")
	}

	userRole := auth.GetUserRole(ctx)

	if userRole != auth.RoleTeacher && userRole != auth.RoleAdmin {
		return nil, status.Error(
			codes.PermissionDenied,
			"user not allowed to delete assignments",
		)
	}

	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid assignment ID format",
		)
	}

	err = s.assignments.DeleteTemplate(ctx, userID, assignmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"assignment not found",
			)
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &assignmentsv1.DeleteAssignmentResponse{}, nil
}

func (s *serverAPI) GetAssignment(
	ctx context.Context,
	req *assignmentsv1.GetAssignmentRequest,
) (*assignmentsv1.GetAssignmentResponse, error) {
	userRole := auth.GetUserRole(ctx)

	if userRole != auth.RoleTeacher && userRole != auth.RoleAdmin {
		return nil, status.Error(
			codes.PermissionDenied,
			"user not allowed to delete assignments",
		)
	}

	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid assignment ID format",
		)
	}

	assignment, targets, err := s.assignments.GetTemplate(ctx, assignmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"assignment with provided ID not found",
			)
		}
		return nil, status.Error(
			codes.Internal,
			"failed to retrieve assignment",
		)
	}

	widgetConfig, err := rawMessageToStructPB(assignment.WidgetConfig)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to convert widget config",
		)
	}

	assignmentProto := &assignmentsv1.AssignmentTemplate{
		Id:           assignment.ID.String(),
		CreatorId:    assignment.CreatorID.String(),
		Title:        assignment.Title,
		Description:  assignment.Description,
		WidgetId:     assignment.WidgetID.String(),
		WidgetConfig: widgetConfig,
		DueDate:      timestamppb.New(*assignment.DueDate),
		CreatedAt:    timestamppb.New(assignment.CreatedAt),
		UpdatedAt:    timestamppb.New(assignment.UpdatedAt),
	}

	targetsProto := make([]*assignmentsv1.AssignmentTarget, 0, len(targets))
	for _, target := range targets {
		var targetProto assignmentsv1.AssignmentTarget
		if target.GroupID != nil {
			targetProto.Target = &assignmentsv1.AssignmentTarget_GroupId{
				GroupId: target.GroupID.String(),
			}
		} else if target.StudentID != nil {
			targetProto.Target = &assignmentsv1.AssignmentTarget_StudentId{
				StudentId: target.StudentID.String(),
			}
		} else {
			continue
		}
		targetsProto = append(targetsProto, &targetProto)
	}

	return &assignmentsv1.GetAssignmentResponse{
		Template: assignmentProto,
		Targets:  targetsProto,
	}, nil
}

func (s *serverAPI) ListAssignments(
	ctx context.Context,
	req *assignmentsv1.ListAssignmentsRequest,
) (*assignmentsv1.ListAssignmentsResponse, error) {
	userIDstr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	targetID, err := uuid.Parse(userIDstr)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user ID in token",
		)
	}

	role := auth.GetUserRole(ctx)

	if role == auth.RoleAdmin && req.CreatorId != "" {
		parsedCreatorID, err := uuid.Parse(req.CreatorId)
		if err != nil {
			return nil, status.Error(
				codes.InvalidArgument,
				"invalid creator ID format",
			)
		}
		targetID = parsedCreatorID
	}

	limit := req.PageSize
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	assignments, token, err := s.assignments.ListTemplates(
		ctx, targetID, int(limit), req.PageToken,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list assignments")
	}

	assignmentsProto := make([]*assignmentsv1.AssignmentTemplateLight, 0, len(assignments))
	for _, item := range assignments {
		itemProto := &assignmentsv1.AssignmentTemplateLight{
			Id:         item.ID.String(),
			Title:      item.Title,
			WidgetType: item.WidgetType,
			DueDate:    timestamppb.New(*item.DueDate),
		}
		assignmentsProto = append(assignmentsProto, itemProto)
	}

	return &assignmentsv1.ListAssignmentsResponse{
		Items:         assignmentsProto,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) ListAssignmentSubmissions(
	ctx context.Context,
	req *assignmentsv1.ListAssignmentSubmissionsRequest,
) (*assignmentsv1.ListAssignmentSubmissionsResponse, error) {
	templateID, err := uuid.Parse(req.TemplateId)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid template ID format",
		)
	}

	limit := req.PageSize
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	submissionStatus := convertProtoStatus(req.GetStatusFilter())

	submissions, token, err := s.feedbacks.ListSubmissions(
		ctx, templateID, int(limit), req.PageToken, submissionStatus,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to list submissions",
		)
	}

	submissionsProto := make([]*assignmentsv1.Submission, 0, len(submissions))
	var versionProto *assignmentsv1.SubmissionVersionLight
	for _, sub := range submissions {
		versionProto = &assignmentsv1.SubmissionVersionLight{
			Id:            sub.LastVersion.ID.String(),
			VersionNumber: *sub.LastVersion.VersionNumber,
			CreatedAt:     timestamppb.New(*sub.LastVersion.CreatedAt),
		}

		itemProto := &assignmentsv1.Submission{
			Id:            sub.ID.String(),
			TemplateId:    sub.TemplateID.String(),
			StudentId:     sub.StudentID.String(),
			Status:        convertSubmissionStatus(sub.Status),
			StartedAt:     timestamppb.New(sub.StartedAt),
			SubmittedAt:   timestamppb.New(*sub.SubmittedAt),
			LatestVersion: versionProto,
		}

		submissionsProto = append(submissionsProto, itemProto)
	}

	return &assignmentsv1.ListAssignmentSubmissionsResponse{
		Items:         submissionsProto,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) GetStudentSubmission(
	ctx context.Context,
	req *assignmentsv1.GetStudentSubmissionRequest,
) (*assignmentsv1.GetStudentSubmissionResponse, error) {
	submissionID, err := uuid.Parse(req.SubmissionId)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid submission ID format",
		)
	}

	details, err := s.feedbacks.GetSubmissionDetails(ctx, submissionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"submission not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve data")
	}

	widgetConfig, err := rawMessageToStructPB(details.Assignment.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse widget config")
	}

	templateProto := &assignmentsv1.AssignmentTemplate{
		Id:           details.Assignment.ID.String(),
		CreatorId:    details.Assignment.CreatorID.String(),
		Title:        details.Assignment.Title,
		Description:  details.Assignment.Description,
		WidgetId:     details.Assignment.WidgetID.String(),
		WidgetConfig: widgetConfig,
		DueDate:      timestamppb.New(*details.Assignment.DueDate),
		CreatedAt:    timestamppb.New(details.Assignment.CreatedAt),
		UpdatedAt:    timestamppb.New(details.Assignment.UpdatedAt),
	}

	var submittedAtPb *timestamppb.Timestamp
	if details.Submission.SubmittedAt != nil {
		submittedAtPb = timestamppb.New(*details.Submission.SubmittedAt)
	}

	submissionProto := &assignmentsv1.Submission{
		Id:          details.Submission.ID.String(),
		TemplateId:  details.Submission.TemplateID.String(),
		StudentId:   details.Submission.StudentID.String(),
		Status:      convertSubmissionStatus(details.Submission.Status),
		StartedAt:   timestamppb.New(details.Submission.StartedAt),
		SubmittedAt: submittedAtPb,
	}

	versionsHistory := make([]*assignmentsv1.SubmissionVersion, 0, len(details.Versions))
	for _, version := range details.Versions {
		versionPayload, err := rawMessageToStructPB(version.Payload)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse version payload")
		}

		versionsHistory = append(versionsHistory, &assignmentsv1.SubmissionVersion{
			Id:               version.ID.String(),
			VersionNumber:    version.VersionNumber,
			Payload:          versionPayload,
			TimeSpentSeconds: durationpb.New(time.Duration(version.TimeSpentSeconds) * time.Second),
			IsAutosave:       version.IsAutosave,
			CreatedAt:        timestamppb.New(version.CreatedAt),
		})
	}

	feedbackHistory := make([]*assignmentsv1.Feedback, 0, len(details.Feedbacks))
	for _, feedback := range details.Feedbacks {
		feedbackPayload, err := rawMessageToStructPB(feedback.Payload)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse feedback payload")
		}

		var textContent string
		if feedback.TextContent != nil {
			textContent = *feedback.TextContent
		}

		feedbackHistory = append(feedbackHistory, &assignmentsv1.Feedback{
			Id:          feedback.ID.String(),
			VersionId:   feedback.VersionID.String(),
			GraderId:    feedback.GraderID.String(),
			TextContent: textContent,
			Payload:     feedbackPayload,
			IsPublished: feedback.IsPublished,
			CreatedAt:   timestamppb.New(feedback.CreatedAt),
		})
	}

	return &assignmentsv1.GetStudentSubmissionResponse{
		Template:   templateProto,
		Submission: submissionProto,
		History:    versionsHistory,
		Feedback:   feedbackHistory,
	}, nil
}

func (s *serverAPI) ProvideFeedback(
	ctx context.Context,
	req *assignmentsv1.ProvideFeedbackRequest,
) (*assignmentsv1.ProvideFeedbackResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid user ID",
		)
	}

	submissionID, err := uuid.Parse(req.GetSubmissionId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid submission ID",
		)
	}

	versionID, err := uuid.Parse(req.GetVersionId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid submission version ID",
		)
	}

	feedbackPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid payload: %v",
			err,
		)
	}

	dto := dto.Feedback{
		VersionID:    versionID,
		SubmissionID: submissionID,
		TextContent:  req.GetTextContent(),
		Payload:      feedbackPayload,
		IsPublished:  req.GetIsPublished(),
	}

	err = s.feedbacks.ProvideFeedback(
		ctx,
		userID,
		dto,
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(
				codes.NotFound,
				"submission version not found",
			)
		}
		return nil, status.Error(
			codes.Internal,
			"failed to provide feedback",
		)
	}

	return &assignmentsv1.ProvideFeedbackResponse{}, nil
}

func (s *serverAPI) ListMyAssignments(
	ctx context.Context,
	req *assignmentsv1.ListMyAssignmentsRequest,
) (*assignmentsv1.ListMyAssignmentsResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid user ID in token",
		)
	}

	limit := req.GetPageSize()
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	statusFilter := convertProtoStatus(req.GetStatusFilter())

	assignments, token, err := s.submissions.ListStudentAssignments(
		ctx,
		userID,
		int(limit),
		req.GetPageToken(),
		statusFilter,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to list assignments",
		)
	}

	items := make([]*assignmentsv1.ListMyAssignmentsResponse_StudentItem, 0, len(assignments))

	for _, item := range assignments {
		itemTemplate := &assignmentsv1.AssignmentTemplateLight{
			Id:         item.AssignmentID.String(),
			Title:      item.Title,
			WidgetType: item.WidgetType,
			DueDate:    timestamppb.New(*item.DueDate),
		}
		itemStatus := convertSubmissionStatus(*item.Status)

		items = append(items, &assignmentsv1.ListMyAssignmentsResponse_StudentItem{
			Template:    itemTemplate,
			Status:      itemStatus,
			HasFeedback: item.HasFeedback,
		})
	}

	return &assignmentsv1.ListMyAssignmentsResponse{
		Items:         items,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) StartAssignment(
	ctx context.Context,
	req *assignmentsv1.StartAssignmentRequest,
) (*assignmentsv1.StartAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid user ID in token",
		)
	}

	templateID, err := uuid.Parse(req.GetTemplateId())
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"invalid template ID",
		)
	}

	submissionID, startedAt, err := s.submissions.StartAssignment(
		ctx,
		userID,
		templateID,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to start assignment",
		)
	}

	return &assignmentsv1.StartAssignmentResponse{
		SubmissionId: submissionID.String(),
		StartedAt:    timestamppb.New(startedAt),
	}, nil
}

func (s *serverAPI) SaveAssignmentDraft(
	ctx context.Context,
	req *assignmentsv1.SaveAssignmentDraftRequest,
) (*assignmentsv1.SaveAssignmentDraftResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid user ID in token",
		)
	}

	submissionPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid payload: %v",
			err,
		)
	}

	submissionID, err := uuid.Parse(req.GetSubmissionId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid submission ID: %v",
			err,
		)
	}

	saveVersion := dto.SaveVersion{
		SubmissionID: submissionID,
		Payload:      submissionPayload,
		TimeSpent:    req.GetTimeSpent().AsDuration(),
	}

	submissionVersionID, err := s.submissions.SaveDraft(
		ctx,
		userID,
		saveVersion,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to save assignment draft",
		)
	}

	return &assignmentsv1.SaveAssignmentDraftResponse{
		VersionId: submissionVersionID.String(),
		SavedAt:   timestamppb.New(time.Now()),
	}, nil
}

func (s *serverAPI) SubmitAssignment(
	ctx context.Context,
	req *assignmentsv1.SubmitAssignmentRequest,
) (*assignmentsv1.SubmitAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"authentication required",
		)
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid user ID in token",
		)
	}

	submissionPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid payload: %v",
			err,
		)
	}

	submissionID, err := uuid.Parse(req.GetSubmissionId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"invalid submission ID: %v",
			err,
		)
	}

	saveVersion := dto.SaveVersion{
		SubmissionID: submissionID,
		Payload:      submissionPayload,
		TimeSpent:    req.GetTimeSpent().AsDuration(),
	}

	submissionID, submissionStatus, err := s.submissions.SubmitAssignment(
		ctx,
		userID,
		saveVersion,
	)
	if err != nil {
		return nil, status.Error(
			codes.Internal,
			"failed to submit assignment",
		)
	}

	return &assignmentsv1.SubmitAssignmentResponse{
		VersionId: submissionID.String(),
		Status:    convertSubmissionStatus(submissionStatus),
	}, nil
}
