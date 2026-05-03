package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	log    *slog.Logger
	server *http.Server
}

func New(log *slog.Logger, addr string, provider JWKSProvider) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /.well-known/jwks.json", newJWKSHandler(provider))

	readHeaderSeconds := 5
	readWriteTimeoutSeconds := 10
	idleTimeoutSeconds := 120

	return &Server{
		log: log,
		server: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: time.Duration(readHeaderSeconds) * time.Second,

			ReadTimeout:  time.Duration(readWriteTimeoutSeconds) * time.Second,
			WriteTimeout: time.Duration(readWriteTimeoutSeconds) * time.Second,
			IdleTimeout:  time.Duration(idleTimeoutSeconds) * time.Second,
		},
	}
}

func (s *Server) Run() error {
	s.log.Info("http server started", slog.String("addr", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.InfoContext(ctx, "stopping http server")
	return s.server.Shutdown(ctx)
}
