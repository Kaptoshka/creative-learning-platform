package middleware

import (
	"context"
	"net/http"
	"strings"

	httputil "github.com/Kaptoshka/creative-learning-platform/gateway/internal/adapters/incoming/http/utils"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
)

type contextKey string

const TokenKey contextKey = "token"

// Auth extracts Bearer token from Authorization header and puts it into context.
// Handlers retrieve it via TokenFromCtx.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			httputil.Error(w, domain.ErrUnauthenticated)
			return
		}

		ctx := context.WithValue(r.Context(), TokenKey, token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// TokenFromCtx retrieves token from context set by Auth middleware.
func TokenFromCtx(ctx context.Context) string {
	token, _ := ctx.Value(TokenKey).(string)
	return token
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}

	return parts[1]
}
