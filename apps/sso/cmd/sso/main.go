package main

import (
	"log/slog"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/app"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/config"
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

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Database.User, cfg.Database.Password),
		Host:   net.JoinHostPort(cfg.Database.Host, strconv.Itoa(cfg.Database.Port)),
		Path:   cfg.Database.Name,
	}

	q := u.Query()
	q.Set("sslmode", cfg.Database.SSLMode)
	u.RawQuery = q.Encode()

	connString := u.String()
	privateKeyPEM := mustLoadSigningKey(cfg.Signing)

	application, err := app.New(
		log,
		cfg.GRPC.Port,
		cfg.HTTP.Address,
		connString,
		cfg.TokenTTL,
		cfg.RefreshTTL,
		privateKeyPEM,
		cfg.Signing.KeyID,
	)
	if err != nil {
		log.Error("failed to initialize application", slog.Any("error", err))
		os.Exit(1)
	}

	go application.GRPCServer.MustRun()

	go func() {
		if err = application.HTTPServer.Run(); err != nil {
			log.Error("http server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sysSign := <-stop

	log.Info("stopping application", slog.String("signal", sysSign.String()))

	application.Stop()

	log.Info("application stopped")
}

func mustLoadSigningKey(cfg config.SigningConfig) []byte {
	if cfg.KeyPEM != "" {
		return []byte(cfg.KeyPEM)
	}

	if cfg.KeyPath == "" {
		panic("signing key not configured: set signing.key_path or $SIGNING_KEY_PEM")
	}

	pem, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		panic("failed to read signing key file: " + err.Error())
	}

	return pem
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
