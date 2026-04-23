package app

import (
	"log/slog"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/adapters/driven/postgres"
	grpcapp "github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/adapters/driving/grpc"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/assignment"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/feedback"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/submission"
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

	assignmentService := assignment.New(log, client.AssignmentRepo, client.SubmissionRepo, client.FeedbackRepo)
	submissionService := submission.New(log, client.AssignmentRepo, client.SubmissionRepo, client.FeedbackRepo)
	feedbackService := feedback.New(log, client.AssignmentRepo, client.SubmissionRepo, client.FeedbackRepo)

	grpcApp := grpcapp.New(log, assignmentService, submissionService, feedbackService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
