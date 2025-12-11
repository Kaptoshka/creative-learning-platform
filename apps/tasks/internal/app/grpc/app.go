package grpcapp

import (
	"log/slog"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}
