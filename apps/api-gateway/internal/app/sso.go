package app

import (
	"context"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type SSOUseCase struct {
	sso outgoing.SSOService
}

func NewSSOUseCase(sso outgoing.SSOService) *SSOUseCase {
	return &SSOUseCase{sso: sso}
}

func (uc *SSOUseCase) Register(ctx context.Context, req domain.RegisterRequest) (domain.RegisterResponse, error) {
	if req.Email == "" {
		return domain.RegisterResponse{}, fmt.Errorf("%w: email is required", domain.ErrInvalidArgument)
	}
	if req.Password == "" {
		return domain.RegisterResponse{}, fmt.Errorf("%w: password is required", domain.ErrInvalidArgument)
	}
	if req.FirstName == "" {
		return domain.RegisterResponse{}, fmt.Errorf("%w: first_name is required", domain.ErrInvalidArgument)
	}
	if req.LastName == "" {
		return domain.RegisterResponse{}, fmt.Errorf("%w: last_name is required", domain.ErrInvalidArgument)
	}

	return uc.sso.Register(ctx, req)
}

func (uc *SSOUseCase) Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
	if req.Email == "" {
		return domain.LoginResponse{}, fmt.Errorf("%w: email is required", domain.ErrInvalidArgument)
	}
	if req.Password == "" {
		return domain.LoginResponse{}, fmt.Errorf("%w: password is required", domain.ErrInvalidArgument)
	}
	if req.AppID == 0 {
		return domain.LoginResponse{}, fmt.Errorf("%w: app_id is required", domain.ErrInvalidArgument)
	}

	return uc.sso.Login(ctx, req)
}

func (uc *SSOUseCase) Logout(ctx context.Context, req domain.LogoutRequest) (domain.LogoutResponse, error) {
	if req.Token == "" {
		return domain.LogoutResponse{}, fmt.Errorf("%w: token is required", domain.ErrInvalidArgument)
	}

	return uc.sso.Logout(ctx, req)
}
