package middleware

import (
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

type responseWriter struct {
	http.ResponseWriter

	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func Logger(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		method := url.QueryEscape(r.Method)
		path := url.PathEscape(r.URL.Path)
		remoteAddr := url.QueryEscape(r.RemoteAddr)

		log.Info(
			"request",
			"method", method,
			"path", path,
			"status", rw.status,
			"duration", time.Since(start).String(),
			"remote_addr", remoteAddr,
		)
	})
}
