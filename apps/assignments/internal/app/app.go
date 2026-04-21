package app

import (
	"log/slog"

	grpcapp "github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/app/grpc"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/service/assignment"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/storage/postgres"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/storage/postgres/submission"
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
