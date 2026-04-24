package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/ports/driving"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	assignmentService driving.AssignmentService,
	submissionService driving.SubmissionService,
	feedbackService driving.FeedbackService,
	port int,
	jwksURL string,
) (*App, error) {
	authInterceptor, err := NewAuthInterceptor(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth interceptor: %w", err)
	}

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authInterceptor.Unary(),
		),
	)

	Register(gRPCServer, assignmentService, submissionService, feedbackService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}, nil
}

// Run runs gRPC server
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running")
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	log.Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
