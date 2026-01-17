package grpcapp

import (
	"log/slog"

	tasksgrpc "tasks/internal/grpc/tasks"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	assignmentService tasksgrpc.Assignments,
	submissionService tasksgrpc.Submissions,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	tasksgrpc.Register(gRPCServer, assignmentService, submissionService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
