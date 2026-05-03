package grpc

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/ports/driving"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new instance of the gRPC app struct.
func New(
	log *slog.Logger,
	authService driving.AuthService,
	port int,
	publicKey *rsa.PublicKey,
) *App {
	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			AdminOnlyInterceptor(publicKey),
		),
	)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthSrv)
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun run gRPC server and panic if error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("failed to run app", "error", err)
		panic(err)
	}
}

// Run runs gRPC and HTTP server.
// TODO: Separate to two functions.
func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	lc := net.ListenConfig{}
	l, err := lc.Listen(context.Background(), "tcp", ":"+strconv.Itoa(a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running")
	if err = a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	a.gRPCServer.GracefulStop()
}
