package app

import (
	"log/slog"
	"time"

	grpcapp "sso/internal/app/grpc"
	"sso/internal/services/auth"
	"sso/internal/storage"
	"sso/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

type AuthUserStorageAdapter struct {
	storage.UserStorage
	storage.RoleStorage
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

	userStorageAdapter := &AuthUserStorageAdapter{
		UserStorage: client.UserStorage,
		RoleStorage: client.RoleStorage,
	}

	authService := auth.New(log, userStorageAdapter, userStorageAdapter, client.AppStorage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
