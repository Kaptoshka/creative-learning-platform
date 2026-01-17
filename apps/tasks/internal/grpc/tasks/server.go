package tasks

import (
	"google.golang.org/grpc"

	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"
)

type Assignments interface {
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
