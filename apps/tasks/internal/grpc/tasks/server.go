package tasks

import (
	"context"
	"errors"

	"tasks/internal/domain/models"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
)

type Assignments interface {
	Update(
		ctx context.Context,
		assignmentID string,
		updates map[string]any,
		targets []*models.AssignmentTarget,
	) (*models.AssignmentTemplate, error)
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

	targets := make([]*models.AssignmentTarget, len(req.Targets))

	for _, trg := range req.Targets {
		target, err := processTarget(trg)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		targets = append(targets, target)
	}

	updateModel, err := s.assignments.Update(ctx, req.AssignmentId, updates, targets)
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
