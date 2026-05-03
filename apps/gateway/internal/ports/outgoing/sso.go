package outgoing

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
)

type SSOService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (domain.RegisterResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error)
	Refresh(ctx context.Context, req domain.RefreshRequest) (domain.RefreshResponse, error)
	Logout(ctx context.Context, req domain.LogoutRequest) (domain.LogoutResponse, error)
	LogoutAll(ctx context.Context, req domain.LogoutAllRequest) (domain.LogoutAllResponse, error)

	// Admin only
	RegisterApp(ctx context.Context, req domain.RegisterAppRequest) (domain.RegisterAppResponse, error)
	DeactivateApp(ctx context.Context, req domain.DeactivateAppRequest) (domain.DeactivateAppResponse, error)
}
