package grpc

import (
	"log/slog"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	assignmentService Assignments,
	submissionService Submissions,
	feedbackService Feedbacks,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	Register(gRPCServer, assignmentService, submissionService, feedbackService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
