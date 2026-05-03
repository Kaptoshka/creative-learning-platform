package app

import (
	"fmt"
	"log/slog"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/adapters/driven/postgres"
	grpcapp "github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/adapters/driving/grpc"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/assignment"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/feedback"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/core/service/submission"
)

type App struct {
	GRPCServer     *grpcapp.App
	postgresClient *postgres.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	jwksURL string,
	connString string,
) (*App, error) {
	client, err := postgres.New(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to storage: %w", err)
	}

	assignmentService := assignment.New(
		log,
		&client.AssignmentRepo,
		&client.SubmissionRepo,
		&client.FeedbackRepo,
	)
	submissionService := submission.New(
		log,
		&client.AssignmentRepo,
		&client.SubmissionRepo,
		&client.FeedbackRepo,
	)
	feedbackService := feedback.New(
		log,
		&client.AssignmentRepo,
		&client.SubmissionRepo,
		&client.FeedbackRepo,
	)

	grpcApp, err := grpcapp.New(
		log,
		assignmentService,
		submissionService,
		feedbackService,
		grpcPort,
		jwksURL,
	)

	return &App{
		GRPCServer:     grpcApp,
		postgresClient: client,
	}, nil
}

func (a *App) Stop() {
	a.GRPCServer.Stop()
	a.postgresClient.Close()
}
