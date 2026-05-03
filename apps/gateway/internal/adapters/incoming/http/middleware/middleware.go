package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/config"
)

type Middleware struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) Public(log *slog.Logger, next http.Handler) http.Handler {
	return Logger(log, CORS(m.cfg)(next))
}

func (m *Middleware) Protected(log *slog.Logger, next http.Handler) http.Handler {
	return Logger(log, CORS(m.cfg)(Auth(log, next)))
}
