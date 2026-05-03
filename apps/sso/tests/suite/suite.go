package suite

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/config"

	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/sso/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T

	Cfg        *config.Config
	AuthClient ssov1.AuthServiceClient
	AppID      string
	JWKSUrl    string
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	configPath := os.Getenv("TEST_CFG_PATH")
	cfg := config.MustLoadByPath(configPath)

	appID := os.Getenv("TEST_APP_ID")
	if appID == "" {
		t.Fatal("TEST_APP_ID environment variable is required")
	}

	jwksURL := fmt.Sprintf(
		"http://%s/.well-known/jwks.json",
		net.JoinHostPort(grpcHost, httpPort(cfg.HTTP.Address)),
	)

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	t.Cleanup(func() {
		_ = cc.Close()
	})

	return ctx, &Suite{
		T:          t,
		Cfg:        cfg,
		AuthClient: ssov1.NewAuthServiceClient(cc),
		AppID:      appID,
		JWKSUrl:    jwksURL,
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}

func httpPort(address string) string {
	_, port, err := net.SplitHostPort(address)
	if err != nil {
		return "9090"
	}
	return port
}
