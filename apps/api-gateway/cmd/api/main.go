package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/app"
	"api-gateway/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("API Gateway started", slog.String("env", cfg.Env))

	application := app.New(log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sysSign := <-stop

	log.Info("stopping API Gateway", slog.String("signal", sysSign.String()))

	log.Info("API Gateway stopped")
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
