package app

import (
	"log/slog"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driven/postgres"
	grpcapp "github.com/Kaptoshka/creative-learning-platform/sso-service/internal/adapters/driving/grpc"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/service/auth"
)

type App struct {
	GRPCServer     *grpcapp.App
	postgresClient *postgres.Storage
}

// New creates a new instance of the App struct.
func New(
	log *slog.Logger,
	grpcPort int,
	connString string,
	tokenTTL time.Duration,
) *App {
	client, err := postgres.New(connString)
	if err != nil {
		return nil
	}

	authService := auth.New(log, &client.AppRepo, &client.RoleRepo, &client.UserRepo, tokenTTL)

	grpcApp := grpcapp.New(log, &authService, grpcPort)

	return &App{
		GRPCServer:     grpcApp,
		postgresClient: client,
	}
}

func (a *App) Stop() {
	a.GRPCServer.Stop()
	a.postgresClient.Close()
}
