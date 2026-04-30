package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/jwt"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres"
	grpcapp "github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driving/grpc"
	httpadapter "github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driving/http"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/service/auth"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/ports/driven"
)

type App struct {
	GRPCServer     *grpcapp.App
	HTTPServer     *httpadapter.Server
	postgresClient *postgres.Storage
	cancelCtx      context.CancelFunc
}

// New creates a new instance of the App struct.
func New(
	log *slog.Logger,
	grpcPort int,
	httpAddr string,
	connString string,
	tokenTTL time.Duration,
	refreshTTL time.Duration,
	privateKey []byte,
	keyID string,
) (*App, error) {
	client, err := postgres.New(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	rsaKey, err := jwt.ParseRSAPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	jwtService := jwt.New(tokenTTL, rsaKey, keyID)

	authService := auth.New(
		log,
		&client.AppRepo,
		&client.RoleRepo,
		&client.UserRepo,
		&client.GroupRepo,
		&client.RefreshRepo,
		jwtService,
		refreshTTL,
	)

	grpcApp := grpcapp.New(log, &authService, grpcPort, jwtService.PublicKey())

	httpServer := httpadapter.New(log, httpAddr, jwtService)

	ctx, cancel := context.WithCancel(context.Background())

	application := &App{
		GRPCServer:     grpcApp,
		HTTPServer:     httpServer,
		postgresClient: client,
		cancelCtx:      cancel,
	}

	go application.runTokenCleanup(ctx, log, &client.RefreshRepo)

	return application, nil
}

func (a *App) runTokenCleanup(
	ctx context.Context,
	log *slog.Logger,
	repo driven.RefreshRepo,
) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := repo.DeleteExpired(ctx); err != nil {
				log.Error("failed to cleanup expired tokens",
					slog.Any("error", err),
				)
			} else {
				log.Info("expired refresh tokens cleaned up")
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *App) Stop() {
	a.cancelCtx()
	a.GRPCServer.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = a.HTTPServer.Stop(ctx)
	a.postgresClient.Close()
}
