package tasks

import (
	"context"
	"errors"
	"time"

	"tasks/internal/domain/models"
	"tasks/internal/storage"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
}

type Submissions interface {
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
