package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/auth"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/dto"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/domain/models"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Assignments interface {
	Create(
		ctx context.Context,
		creatorID uuid.UUID,
		dto *dto.CreateAssignment,
	) (uuid.UUID, error)
	Update(
		ctx context.Context,
		assignmentID uuid.UUID,
		updates map[string]any,
		targets []*dto.Target,
	) (*models.AssignmentTemplate, error)
	Delete(
		ctx context.Context,
		assignmentID uuid.UUID,
	) error
	GetByID(
		ctx context.Context,
		assignmentID uuid.UUID,
	) (*models.AssignmentTemplate, []*dto.Target, error)
	List(
		ctx context.Context,
		creatorID uuid.UUID,
		pageSize int32,
		pageToken string,
	) ([]*models.AssignmentTemplateLight, string, error)
	ListByUserID(
		ctx context.Context,
		userID uuid.UUID,
		pageSize int32,
		pageToken string,
		statusFilter domain.SubmissionStatus,
	) ([]*models.AssignmentTemplateLight, string, error)
	Start(
		ctx context.Context,
		studentID uuid.UUID,
		templateID uuid.UUID,
	) (uuid.UUID, time.Time, error)
	SaveDraft(
		ctx context.Context,
		studentID uuid.UUID,
		payload json.RawMessage,
	) (uuid.UUID, error)
	Submit(
		ctx context.Context,
		studentID uuid.UUID,
		payload json.RawMessage,
	) (uuid.UUID, domain.SubmissionStatus, error)
}

type Submissions interface {
	ListByTemplateID(
		ctx context.Context,
		templateID uuid.UUID,
		pageSize int32,
		pageToken string,
	) ([]*models.SubmissionItem, string, error)
	GetByID(
		ctx context.Context,
		submissionID uuid.UUID,
	) (
		*models.AssignmentTemplate,
		*models.Submission,
		[]*models.SubmissionVersion,
		[]*models.Feedback,
		error,
	)
}

type Feedbacks interface {
	Provide(
		ctx context.Context,
		graderID uuid.UUID,
		dto *dto.Feedback,
	) error
}

type serverAPI struct {
	tasksv1.UnimplementedTasksServiceServer
	assignments Assignments
	submissions Submissions
	feedbacks   Feedbacks
}

func Register(
	gRPC *grpc.Server,
	assignments Assignments,
	submissions Submissions,
) {
	tasksv1.RegisterTasksServiceServer(gRPC, &serverAPI{
		assignments: assignments,
		submissions: submissions,
	})
}

func (s *serverAPI) CreateAssignment(
	ctx context.Context,
	req *tasksv1.CreateAssignmentRequest,
) (*tasksv1.CreateAssignmentResponse, error) {
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

		targets = append(targets, target)
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

	assignmentDTO := &dto.CreateAssignment{
		Title:        req.Title,
		Description:  req.Description,
		WidgetID:     widgetID,
		WidgetConfig: widgetConfig,
		DueDate:      &dueTime,
		Targets:      targets,
	}

	assignmentID, err := s.assignments.Create(
		ctx,
		userID,
		assignmentDTO,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "creating assignment error")
	}

	return &tasksv1.CreateAssignmentResponse{
		Id: assignmentID.String(),
	}, nil
}

func (s *serverAPI) UpdateAssignment(
	ctx context.Context,
	req *tasksv1.UpdateAssignmentRequest,
) (*tasksv1.UpdateAssignmentResponse, error) {
	userRole := auth.GetUserRole(ctx)
	if userRole != auth.RoleAdmin && userRole != auth.RoleTeacher {
		return nil, status.Error(codes.PermissionDenied, "user not allowed to update assignment")
	}

	if !req.UpdateMask.IsValid(req.Template) {
		return nil, status.Error(codes.InvalidArgument, "invalid update mask")
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

	targets := make([]*dto.Target, 0, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid target")
		}

		targets = append(targets, target)
	}

	assignmentID, err := uuid.Parse(req.AssignmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "assignment ID cannot be parsed")
	}

	updateModel, err := s.assignments.Update(ctx, assignmentID, updates, targets)
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot update assignment")
	}

	widgetConfig, err := rawMessageToStructPB(updateModel.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tasksv1.UpdateAssignmentResponse{
		Template: &tasksv1.AssignmentTemplate{
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
	t *tasksv1.AssignmentTarget,
) (*dto.Target, error) {
	switch v := t.GetTarget().(type) {
	case *tasksv1.AssignmentTarget_GroupId:
		groupID, err := uuid.Parse(v.GroupId)
		if err != nil {
			return nil, errors.New("invalid group ID")
		}

		return &dto.Target{
			GroupID: &groupID,
		}, nil
	case *tasksv1.AssignmentTarget_StudentId:
		studentID, err := uuid.Parse(v.StudentId)
		if err != nil {
			return nil, errors.New("invalid student ID")
		}

		return &dto.Target{
			StudentID: &studentID,
		}, nil

	case nil:
		return nil, nil

	default:
		return nil, errors.New("unknown target type")
	}
}

func (s *serverAPI) DeleteAssignment(
	ctx context.Context,
	req *tasksv1.DeleteAssignmentRequest,
) (*tasksv1.DeleteAssignmentResponse, error) {
	userRole := auth.GetUserRole(ctx)

	if userRole != auth.RoleTeacher && userRole != auth.RoleAdmin {
		return nil, status.Error(codes.PermissionDenied, "user not allowed to delete assignments")
	}

	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assignment ID format")
	}

	err = s.assignments.Delete(ctx, assignmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "assignment not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tasksv1.DeleteAssignmentResponse{}, nil
}

func (s *serverAPI) GetAssignment(
	ctx context.Context,
	req *tasksv1.GetAssignmentRequest,
) (*tasksv1.GetAssignmentResponse, error) {
	userRole := auth.GetUserRole(ctx)

	if userRole != auth.RoleTeacher && userRole != auth.RoleAdmin {
		return nil, status.Error(codes.PermissionDenied, "user not allowed to delete assignments")
	}

	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assignment ID format")
	}

	assignment, targets, err := s.assignments.GetByID(ctx, assignmentID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "assignment with provided ID not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve assignment")
	}

	widgetConfig, err := rawMessageToStructPB(assignment.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert widget config")
	}

	assignmentProto := &tasksv1.AssignmentTemplate{
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

	targetsProto := make([]*tasksv1.AssignmentTarget, 0, len(targets))
	for _, target := range targets {
		var targetProto *tasksv1.AssignmentTarget
		if target.GroupID != nil {
			targetProto.Target = &tasksv1.AssignmentTarget_GroupId{
				GroupId: target.GroupID.String(),
			}
		} else if target.StudentID != nil {
			targetProto.Target = &tasksv1.AssignmentTarget_StudentId{
				StudentId: target.StudentID.String(),
			}
		} else {
			continue
		}
		targetsProto = append(targetsProto, targetProto)
	}

	return &tasksv1.GetAssignmentResponse{
		Template: assignmentProto,
		Targets:  targetsProto,
	}, nil
}

func (s *serverAPI) ListAssignments(
	ctx context.Context,
	req *tasksv1.ListAssignmentsRequest,
) (*tasksv1.ListAssignmentsResponse, error) {
	userIDstr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	targetID, err := uuid.Parse(userIDstr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID in token")
	}

	role := auth.GetUserRole(ctx)

	if role == auth.RoleAdmin && req.CreatorId != "" {
		parsedCreatorID, err := uuid.Parse(req.CreatorId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid creator ID format")
		}
		targetID = parsedCreatorID
	}

	limit := req.PageSize
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	assignments, token, err := s.assignments.List(ctx, targetID, limit, req.PageToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list assignments")
	}

	assignmentsProto := make([]*tasksv1.AssignmentTemplateLight, 0, len(assignments))
	for _, item := range assignments {
		itemProto := &tasksv1.AssignmentTemplateLight{
			Id:         item.ID.String(),
			Title:      item.Title,
			WidgetType: item.WidgetType,
			DueDate:    timestamppb.New(*item.DueDate),
		}
		assignmentsProto = append(assignmentsProto, itemProto)
	}

	return &tasksv1.ListAssignmentsResponse{
		Items:         assignmentsProto,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) ListAssignmentSubmissions(
	ctx context.Context,
	req *tasksv1.ListAssignmentSubmissionsRequest,
) (*tasksv1.ListAssignmentSubmissionsResponse, error) {
	templateID, err := uuid.Parse(req.TemplateId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid template ID format")
	}

	limit := req.PageSize
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	submissions, token, err := s.submissions.ListByTemplateID(ctx, templateID, limit, req.PageToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list submissions")
	}

	submissionsProto := make([]*tasksv1.Submission, 0, len(submissions))
	for _, item := range submissions {
		payload, err := rawMessageToStructPB(item.SubmissionVersion.Payload)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse submission payload")
		}

		submissionVersion := &tasksv1.SubmissionVersion{
			Id:               item.SubmissionVersion.ID.String(),
			VersionNumber:    item.SubmissionVersion.VersionNumber,
			Payload:          payload,
			TimeSpentSeconds: durationpb.New(item.SubmissionVersion.TimeSpentSeconds),
			IsAutosave:       item.SubmissionVersion.IsAutosave,
			CreatedAt:        timestamppb.New(item.SubmissionVersion.CreatedAt),
		}

		itemProto := &tasksv1.Submission{
			Id:            item.Submission.ID.String(),
			TemplateId:    item.Submission.TemplateID.String(),
			StudentId:     item.Submission.StudentID.String(),
			Status:        convertSubmissionStatus(item.Submission.Status),
			StartedAt:     timestamppb.New(item.Submission.StartedAt),
			SubmittedAt:   timestamppb.New(*item.Submission.SubmittedAt),
			LatestVersion: submissionVersion,
		}

		submissionsProto = append(submissionsProto, itemProto)
	}

	return &tasksv1.ListAssignmentSubmissionsResponse{
		Items:         submissionsProto,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) GetStudentSubmission(
	ctx context.Context,
	req *tasksv1.GetStudentSubmissionRequest,
) (*tasksv1.GetStudentSubmissionResponse, error) {
	submissionID, err := uuid.Parse(req.SubmissionId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid submission ID format")
	}

	assignment, submission, submissionVersions, feedbacks, err := s.submissions.GetByID(ctx, submissionID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "submission not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve data")
	}

	widgetConfig, err := rawMessageToStructPB(assignment.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse widget config")
	}

	templateProto := &tasksv1.AssignmentTemplate{
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

	var submittedAtPb *timestamppb.Timestamp
	if submission.SubmittedAt != nil {
		submittedAtPb = timestamppb.New(*submission.SubmittedAt)
	}

	submissionProto := &tasksv1.Submission{
		Id:          submission.ID.String(),
		TemplateId:  submission.TemplateID.String(),
		StudentId:   submission.StudentID.String(),
		Status:      convertSubmissionStatus(submission.Status),
		StartedAt:   timestamppb.New(submission.StartedAt),
		SubmittedAt: submittedAtPb,
	}

	versionsHistory := make([]*tasksv1.SubmissionVersion, 0, len(submissionVersions))
	for _, version := range submissionVersions {
		versionPayload, err := rawMessageToStructPB(version.Payload)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse version payload")
		}

		versionsHistory = append(versionsHistory, &tasksv1.SubmissionVersion{
			Id:               version.ID.String(),
			VersionNumber:    version.VersionNumber,
			Payload:          versionPayload,
			TimeSpentSeconds: durationpb.New(time.Duration(version.TimeSpentSeconds) * time.Second),
			IsAutosave:       version.IsAutosave,
			CreatedAt:        timestamppb.New(version.CreatedAt),
		})
	}

	feedbackHistory := make([]*tasksv1.Feedback, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		feedbackPayload, err := rawMessageToStructPB(feedback.Payload)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to parse feedback payload")
		}

		var textContent string
		if feedback.TextContent != nil {
			textContent = *feedback.TextContent
		}

		feedbackHistory = append(feedbackHistory, &tasksv1.Feedback{
			Id:          feedback.ID.String(),
			VersionId:   feedback.VersionID.String(),
			GraderId:    feedback.GraderID.String(),
			TextContent: textContent,
			Payload:     feedbackPayload,
			IsPublished: feedback.IsPublished,
			CreatedAt:   timestamppb.New(feedback.CreatedAt),
		})
	}

	return &tasksv1.GetStudentSubmissionResponse{
		Template:   templateProto,
		Submission: submissionProto,
		History:    versionsHistory,
		Feedback:   feedbackHistory,
	}, nil
}

func (s *serverAPI) ProvideFeedback(
	ctx context.Context,
	req *tasksv1.ProvideFeedbackRequest,
) (*tasksv1.ProvideFeedbackResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	submissionID, err := uuid.Parse(req.GetSubmissionId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid submission ID")
	}

	versionID, err := uuid.Parse(req.GetVersionId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid submission version ID")
	}

	feedbackPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload: %v", err)
	}

	dto := &dto.Feedback{
		VersionID:    versionID,
		SubmissionID: submissionID,
		TextContent:  req.GetTextContent(),
		Payload:      feedbackPayload,
		IsPublished:  req.GetIsPublished(),
	}

	err = s.feedbacks.Provide(
		ctx,
		userID,
		dto,
	)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "submission version not found")
		}
		return nil, status.Error(codes.Internal, "failed to provide feedback")
	}

	return &tasksv1.ProvideFeedbackResponse{}, nil
}

func (s *serverAPI) ListMyAssignments(
	ctx context.Context,
	req *tasksv1.ListMyAssignmentsRequest,
) (*tasksv1.ListMyAssignmentsResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
	}

	limit := req.GetPageSize()
	if limit <= 0 {
		limit = domain.DefaultPageSizeLimit
	}

	statusFilter := convertProtoStatus(req.GetStatusFilter())

	assignments, token, err := s.assignments.ListByUserID(
		ctx,
		userID,
		limit,
		req.GetPageToken(),
		statusFilter,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list assignments")
	}

	items := make([]*tasksv1.ListMyAssignmentsResponse_StudentItem, 0, len(assignments))

	for _, item := range assignments {
		itemTemplate := &tasksv1.AssignmentTemplateLight{
			Id:         item.ID.String(),
			Title:      item.Title,
			WidgetType: item.WidgetType,
			DueDate:    timestamppb.New(*item.DueDate),
		}
		itemStatus := convertSubmissionStatus(item.Status)
		hasFeedback := false
		if item.Status == domain.StatusGraded || item.Status == domain.StatusReturned {
			hasFeedback = true
		}

		items = append(items, &tasksv1.ListMyAssignmentsResponse_StudentItem{
			Template:    itemTemplate,
			Status:      itemStatus,
			HasFeedback: hasFeedback,
		})
	}

	return &tasksv1.ListMyAssignmentsResponse{
		Items:         items,
		NextPageToken: token,
	}, nil
}

func (s *serverAPI) StartAssignment(
	ctx context.Context,
	req *tasksv1.StartAssignmentRequest,
) (*tasksv1.StartAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
	}

	templateID, err := uuid.Parse(req.GetTemplateId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid template ID")
	}

	submissionID, startedAt, err := s.assignments.Start(
		ctx,
		userID,
		templateID,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to start assignment")
	}

	return &tasksv1.StartAssignmentResponse{
		SubmissionId: submissionID.String(),
		StartedAt:    timestamppb.New(startedAt),
	}, nil
}

func (s *serverAPI) SaveAssignmentDraft(
	ctx context.Context,
	req *tasksv1.SaveAssignmentDraftRequest,
) (*tasksv1.SaveAssignmentDraftResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
	}

	submissionPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload: %v", err)
	}

	submissionVersionID, err := s.assignments.SaveDraft(
		ctx,
		userID,
		submissionPayload,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save assignment draft")
	}

	return &tasksv1.SaveAssignmentDraftResponse{
		VersionId: submissionVersionID.String(),
		SavedAt:   timestamppb.New(time.Now()),
	}, nil
}

func (s *serverAPI) SubmitAssignment(
	ctx context.Context,
	req *tasksv1.SubmitAssignmentRequest,
) (*tasksv1.SubmitAssignmentResponse, error) {
	userIDStr, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid user ID in token")
	}

	submissionPayload, err := structPBToRawMessage(req.GetPayload())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload: %v", err)
	}

	submissionID, submissionStatus, err := s.assignments.Submit(
		ctx,
		userID,
		submissionPayload,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to submit assignment")
	}

	return &tasksv1.SubmitAssignmentResponse{
		VersionId: submissionID.String(),
		Status:    convertSubmissionStatus(submissionStatus),
	}, nil
}
