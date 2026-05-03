package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/app"
	"github.com/Kaptoshka/creative-learning-platform/assignment-service/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application")

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	application, err := app.New(log, cfg.GRPC.Port, cfg.GRPC.JWKSURL, connString)
	if err != nil {
		log.Error("failed to initialize application", slog.Any("error", err))
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		if err := application.GRPCServer.Run(); err != nil {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sysSign := <-stop:
		log.Info("stopping application", slog.String("signal", sysSign.String()))
	case err := <-errCh:
		log.Error("gRPC server failed", slog.Any("error", err))
	}

	application.Stop()
	log.Info("application stopped")
}

// setupLogger creates a new logger instance based on the environment.
func setupLogger(env string) *slog.Logger {
	switch env {
	case envLocal:
		return slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
}
