package tasks

import (
	"context"
	"errors"
	"time"

	"tasks/internal/auth"
	"tasks/internal/domain/models"
	"tasks/internal/storage"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Assignments interface {
	Create(
		ctx context.Context,
		title string,
		description string,
		widgetID uuid.UUID,
		widgetConfig *models.JSONB,
		dueDate *time.Time,
		targets []*models.AssignmentTarget,
	) (uuid.UUID, error)
	Update(
		ctx context.Context,
		assignmentID uuid.UUID,
		updates map[string]any,
		targets []*models.AssignmentTarget,
	) (*models.AssignmentTemplate, error)
	Delete(
		ctx context.Context,
		assignmentID uuid.UUID,
	) error
	GetByID(
		ctx context.Context,
		assignmentID uuid.UUID,
	) (*models.AssignmentTemplate, []*models.AssignmentTarget, error)
	List(
		ctx context.Context,
		creatorID uuid.UUID,
		pageSize int32,
		pageToken string,
	) ([]*models.AssignmentTemplateLight, string, error)
}

type Submissions interface {
	ListByTemplateID(
		ctx context.Context,
		templateID uuid.UUID,
		pageSize int32,
		pageToken string,
	) ([]*models.SubmissionItem, string, error)
	// TODO: make simpler output
	GetByID(
		ctx context.Context,
		submissionID uuid.UUID,
	) (*models.AssignmentTemplate, *models.Submission, []*models.SubmissionVersion, []*models.Feedback, error)
}

type serverAPI struct {
	tasksv1.UnimplementedTasksServer
	assignments Assignments
	submissions Submissions
}

func Register(
	gRPC *grpc.Server,
	assignments Assignments,
	submissions Submissions,
) {
	tasksv1.RegisterTasksServer(gRPC, &serverAPI{
		assignments: assignments,
		submissions: submissions,
	})
}

func (s *serverAPI) CreateAssignment(
	ctx context.Context,
	req *tasksv1.CreateAssignmentRequest,
) (*tasksv1.CreateAssignmentResponse, error) {
	targets := make([]*models.AssignmentTarget, 0, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		targets = append(targets, target)
	}

	widgetID, err := uuid.Parse(req.WidgetId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "parsing widget ID error")
	}

	widgetConfig := models.JSONB(req.WidgetConfig.AsMap())

	dueTime := req.DueDate.AsTime()

	assignmentID, err := s.assignments.Create(
		ctx,
		req.Title,
		req.Description,
		widgetID,
		&widgetConfig,
		&dueTime,
		targets,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tasksv1.CreateAssignmentResponse{
		Id: assignmentID.String(),
	}, nil
}

func (s *serverAPI) UpdateAssignment(
	ctx context.Context,
	req *tasksv1.UpdateAssignmentRequest,
) (*tasksv1.AssignmentTemplate, error) {
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

	targets := make([]*models.AssignmentTarget, 0, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		targets = append(targets, target)
	}

	assignmentID, err := uuid.Parse(req.AssignmentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "assignment ID cannot be parsed")
	}

	updateModel, err := s.assignments.Update(ctx, assignmentID, updates, targets)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	widgetConfig, err := structpb.NewStruct(updateModel.WidgetConfig)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &tasksv1.AssignmentTemplate{
		Id:           updateModel.ID.String(),
		CreatorId:    updateModel.CreatorID.String(),
		Title:        updateModel.Title,
		Description:  updateModel.Description,
		WidgetId:     updateModel.WidgetID.String(),
		WidgetConfig: widgetConfig,
		DueDate:      timestamppb.New(updateModel.DueDate),
		CreatedAt:    timestamppb.New(updateModel.CreatedAt),
		UpdatedAt:    timestamppb.New(updateModel.UpdatedAt),
	}, nil
}

func processTarget(t *tasksv1.AssignmentTarget) (*models.AssignmentTarget, error) {
	switch v := t.GetTarget().(type) {
	case *tasksv1.AssignmentTarget_GroupId:
		groupID, err := uuid.Parse(v.GroupId)
		if err != nil {
			return nil, errors.New("invalid group ID")
		}

		return &models.AssignmentTarget{
			GroupID: &groupID,
		}, nil
	case *tasksv1.AssignmentTarget_StudentId:
		studentID, err := uuid.Parse(v.StudentId)
		if err != nil {
			return nil, errors.New("invalid student ID")
		}

		return &models.AssignmentTarget{
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
) (*emptypb.Empty, error) {
	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assignment ID format")
	}

	err = s.assignments.Delete(ctx, assignmentID)
	if err != nil {
		if errors.Is(err, storage.ErrAssignmentNotFound) {
			return nil, status.Error(codes.NotFound, "assignment not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverAPI) GetAssignment(
	ctx context.Context,
	req *tasksv1.GetAssignmentRequest,
) (*tasksv1.GetAssignmentResponse, error) {
	assignmentID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid assignment ID format")
	}

	assignment, targets, err := s.assignments.GetByID(ctx, assignmentID)
	if err != nil {
		if errors.Is(err, storage.ErrAssignmentNotFound) {
			return nil, status.Error(codes.NotFound, "assignment with provided ID not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve assignment")
	}

	widgetConfig, err := structpb.NewStruct(assignment.WidgetConfig)
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
		DueDate:      timestamppb.New(assignment.DueDate),
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
		limit = models.DefaultPageSizeLimit
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
			DueDate:    timestamppb.New(item.DueDate),
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
		limit = models.DefaultPageSizeLimit
	}

	submissions, token, err := s.submissions.ListByTemplateID(ctx, templateID, limit, req.PageToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list submissions")
	}

	submissionsProto := make([]*tasksv1.Submission, 0, len(submissions))
	for _, item := range submissions {
		payload, err := structpb.NewStruct(item.SubmissionVersion.Payload)
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
			StartedAt:     timestamppb.New(item.Submission.StartedAt),
			SubmittedAt:   timestamppb.New(item.Submission.SubmittedAt),
			LatestVersion: submissionVersion,
		}

		switch item.Submission.Status {
		case models.StatusNotSpecified:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_UNSPECIFIED
		case models.StatusNotStarted:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_NOT_STARTED
		case models.StatusInProgress:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_IN_PROGRESS
		case models.StatusSubmitted:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_SUBMITTED
		case models.StatusGraded:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_GRADED
		case models.StatusReturned:
			itemProto.Status = tasksv1.SubmissionStatus_SUBMISSION_STATUS_RETURNED
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
		if errors.Is(err, storage.ErrSubmissionNotFound) {
			return nil, status.Error(codes.NotFound, "submission not found")
		}
		return nil, status.Error(codes.Internal, "failed to retrieve data")
	}

	widgetConfig, err := structpb.NewStruct(assignment.WidgetConfig)
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
		DueDate:      timestamppb.New(assignment.DueDate),
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
		versionPayload, err := structpb.NewStruct(version.Payload)
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
		feedbackPayload, err := structpb.NewStruct(feedback.Payload)
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

func convertSubmissionStatus(status models.SubmissionStatus) tasksv1.SubmissionStatus {
	switch status {
	case models.StatusNotStarted:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_NOT_STARTED
	case models.StatusInProgress:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_IN_PROGRESS
	case models.StatusSubmitted:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_SUBMITTED
	case models.StatusGraded:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_GRADED
	case models.StatusReturned:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_RETURNED
	default:
		return tasksv1.SubmissionStatus_SUBMISSION_STATUS_UNSPECIFIED
	}
}
