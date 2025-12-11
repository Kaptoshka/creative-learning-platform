package grpcconn

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"api-gateway/internal/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
)

type GRPCClient[T any] struct {
	API  T
	Log  *slog.Logger
	Conn *grpc.ClientConn
}

// Close closes the connection.
func (c *GRPCClient[T]) Close() error {
	return c.Conn.Close()
}

func New[T any](
	ctx context.Context,
	log *slog.Logger,
	cfg *config.GRPCClient,
	maker func(grpc.ClientConnInterface) T,
) (*GRPCClient[T], error) {
	const op = "grpcconn.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(cfg.RetriesCount)),
		grpcretry.WithBackoff(grpcretry.BackoffExponential(time.Second)),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadSent, grpclog.PayloadReceived),
	}

	var credentialsOpts grpc.DialOption

	if cfg.Insecure {
		credentialsOpts = grpc.WithTransportCredentials(insecure.NewCredentials())
	}

	cc, err := grpc.NewClient(cfg.Address,
		credentialsOpts,
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	clientAPI := maker(cc)

	return &GRPCClient[T]{
		API:  clientAPI,
		Log:  log,
		Conn: cc,
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger.
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
