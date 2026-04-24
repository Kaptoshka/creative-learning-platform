package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/config"
	assignmentsv1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/proto/assignments/v1"
	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/proto/sso/v1"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/handlers"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/router"
	grpcadapters "github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/outgoing/grpc"
	usecase "github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/use_case"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log        *slog.Logger
	HTTPServer *http.Server
	timeout    time.Duration
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) (*App, error) {
	ssoConn, err := mustDial(log, cfg.Clients.SSO)
	if err != nil {
		return nil, fmt.Errorf("failed to dial sso server: %w", err)
	}
	defer ssoConn.Close()

	assignmentsConn, err := mustDial(log, cfg.Clients.Assignments)
	if err != nil {
		return nil, fmt.Errorf("failed to dial assignments server: %w", err)
	}
	defer assignmentsConn.Close()

	ssoClient := ssov1.NewAuthServiceClient(ssoConn)
	assignmentsClient := assignmentsv1.NewAssignmentServiceClient(assignmentsConn)

	ssoAdapter := grpcadapters.NewSSOAdapter(ssoClient)
	assignmentsAdapter := grpcadapters.NewAssignmentsAdapter(assignmentsClient)

	ssoUseCase := usecase.NewSSOUseCase(ssoAdapter)
	assignmentsUseCase := usecase.NewAssignmentsUseCase(assignmentsAdapter)

	ssoHandler := handlers.NewSSOHandler(ssoUseCase)
	assignmentsHandler := handlers.NewAssignmentsHandler(assignmentsUseCase)

	mw := middleware.New(cfg)

	r := router.New(ssoHandler, assignmentsHandler, mw)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout * 3,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	return &App{
		log:        log,
		HTTPServer: srv,
		timeout:    cfg.HTTPServer.Timeout,
	}, nil
}

func mustDial(
	log *slog.Logger,
	client config.GRPCClient,
) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if !client.Insecure {
		log.Warn("insecure=false but TLS not configured, falling back to insecure", "addr", client.Address)
	}

	conn, err := grpc.NewClient(client.Address, opts...)
	if err != nil {
		slog.Error("failed to connect to gRPC service", "addr", client.Address, "err", err)
		return nil, fmt.Errorf("failed to connect to gRPC service: %w", err)
	}

	log.Info("gRPC client created", "addr", client.Address)
	return conn, nil
}

func (a *App) Stop() {
	const op = "app.Stop"

	log := a.log.With(
		slog.String("op", op),
	)

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	if err := a.HTTPServer.Shutdown(ctx); err != nil {
		log.Error("forced shutdown", "error", err)
	}
}
