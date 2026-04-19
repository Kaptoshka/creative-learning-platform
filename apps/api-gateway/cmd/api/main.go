package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/sso/v1"
	tasksv1 "github.com/Kaptoshka/creative-learning-platform/libs/gen/go/tasks/v1"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/app"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/config"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/handlers"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/middleware"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/router"
	grpcadapters "github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/outgoing/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	ssoConn := mustDial(cfg.Clients.SSO)
	defer ssoConn.Close()

	assignmentsConn := mustDial(cfg.Clients.Assignments)
	defer AssignmentsConn.Close()

	ssoClient := ssov1.NewAuthServiceClient(ssoConn)
	assignmentsClient := tasksv1.NewAssignmentsServiceClient(assignmentsConn)

	ssoAdapter := grpcadapters.NewSSOAdapter(ssoClient)
	assignmentsAdapter := grpcadapters.NewAssignmentsAdapter(assignmentsClient)

	ssoUseCase := app.NewSSOUseCase(ssoAdapter)
	assignmentsUseCase := app.NewAssignmentsUseCase(assignmentsAdapter)

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("gateway starting", "addr", cfg.HTTPServer.Address, "env", cfg.Env, "version", cfg.Version)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPServer.Timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("forced shutdown", "err", err)
	}

	slog.Info("gateway stopped")
}

func mustDial(client config.GRPCClient) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if !client.Insecure {
		slog.Warn("insecure=false but TLS not configured, falling back to insecure", "addr", client.Address)
	}

	conn, err := grpc.NewClient(client.Address, opts...)
	if err != nil {
		slog.Error("failed to connect to gRPC service", "addr", client.Address, "err", err)
		os.Exit(1)
	}

	slog.Info("gRPC client created", "addr", client.Address)
	return conn
}

// setupLogger creates a new logger instance based on the environment.
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
