package app

import (
	"log/slog"

	grpcapp "tasks/internal/app/grpc"
	"tasks/internal/services/assignment"
	"tasks/internal/services/submission"
	"tasks/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	connString string,
) *App {
	client, err := postgres.New(connString)
	if err != nil {
		return nil
	}

	assignmentService := assignment.New(log, client.AssignmentStorage)
	submissionService := submission.New(log, client.SubmissionStorage)

	grpcApp := grpcapp.New(log, assignmentService, submissionService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
