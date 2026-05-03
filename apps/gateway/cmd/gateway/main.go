package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kaptoshka/creative-learning-platform/gateway/app"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting api gateway", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	application, err := app.New(log, cfg)
	if err != nil {
		log.Error("failed to init app", slog.Any("error", err))
		os.Exit(1)
	}

	errCh := make(chan error, 1)
	go func() {
		log.Info("gateway starting", "addr", cfg.HTTPServer.Address, "env", cfg.Env, "version", cfg.Version)
		if err = application.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sysSign := <-stop:
		log.Info("stopping api gateway...", slog.String("signal", sysSign.String()))
	case err = <-errCh:
		log.Error("api gateway server failed", slog.Any("error", err))
	}

	application.Stop()
	log.Info("gateway stopped")
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
