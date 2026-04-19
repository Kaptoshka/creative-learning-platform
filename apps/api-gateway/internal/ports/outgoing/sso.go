package outgoing

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/services/gateway/internal/domain"
)

type SSOService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (domain.RegisterResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error)
	Logout(ctx context.Context, req domain.LogoutRequest) (domain.LogoutResponse, error)
}
