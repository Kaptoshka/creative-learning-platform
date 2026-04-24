package middleware

import (
	"net/http"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/config"
)

type Middleware struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) Public(next http.Handler) http.Handler {
	return Logger(CORS(m.cfg)(next))
}

func (m *Middleware) Protected(next http.Handler) http.Handler {
	return Logger(CORS(m.cfg)(Auth(next)))
}
