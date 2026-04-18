package grpcapp

import (
	"log/slog"

	tasksgrpc "github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/transport/grpc/server"

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
