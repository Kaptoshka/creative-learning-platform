package auth

import (
	"context"

	"tasks/internal/domain/models"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
)

type Tasks interface {
	SubmissionByAssignmentID(
		ctx context.Context,
		assignmentID int64,
	) ([]models.Assignment, error)
	SubmissionByAssignmentIDAndStudentID(
		ctx context.Context,
		assignmentID int64,
		studentID int64,
	) (*models.Submission, error)
	AssignmentByID(
		ctx context.Context,
		assignmentID int64,
	) (*models.Assignment, error)
}

type serverAPI struct {
	tasksv1.UnimplementedTasksServer
	tasks Tasks
}

func Register(gRPC *grpc.Server, tasks Tasks) {
	tasksv1.RegisterTasksServer(gRPC, &serverAPI{tasks: tasks})
}

func (s *serverAPI) Assignment(
	ctx context.Context,
	req *tasksv1.GetAssignmentRequest,
) (*tasksv1.GetAssignmentResponse, error) {
	if err := validateAssignment(req); err != nil {
		return nil, err
	}

	assignment, err := s.tasks.AssignmentByID(
		ctx,
		req.GetAssignmentId(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoContent, err := structpb.NewStruct(assignment.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoDeadline := timestamppb.New(assignment.Deadline)

	return &tasksv1.GetAssignmentResponse{
		Assignment: &tasksv1.Assignment{
			Id:        assignment.ID,
			TeacherId: assignment.TeacherID,
			Title:     assignment.Title,
			Content:   protoContent,
			Deadline:  protoDeadline,
		},
	}, nil
}

func validateAssignment(req *tasksv1.GetAssignmentRequest) error {
	if req.GetAssignmentId() == 0 {
		return status.Error(codes.InvalidArgument, "assignment_id is required")
	}
	return nil
}
