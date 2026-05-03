package usecase

import (
	"context"
	"fmt"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/ports/outgoing"
)

type SSOUseCase struct {
	sso outgoing.SSOService
}

func NewSSOUseCase(sso outgoing.SSOService) *SSOUseCase {
	return &SSOUseCase{sso: sso}
}

func (uc *SSOUseCase) Register(
	ctx context.Context,
	req domain.RegisterRequest,
) (domain.RegisterResponse, error) {
	if req.Email == "" {
		return domain.RegisterResponse{}, fmt.Errorf(
			"%w: email is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.Password == "" {
		return domain.RegisterResponse{}, fmt.Errorf(
			"%w: password is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.FirstName == "" {
		return domain.RegisterResponse{}, fmt.Errorf(
			"%w: first_name is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.LastName == "" {
		return domain.RegisterResponse{}, fmt.Errorf(
			"%w: last_name is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.Register(ctx, req)
}

func (uc *SSOUseCase) Login(
	ctx context.Context,
	req domain.LoginRequest,
) (domain.LoginResponse, error) {
	if req.Email == "" {
		return domain.LoginResponse{}, fmt.Errorf(
			"%w: email is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.Password == "" {
		return domain.LoginResponse{}, fmt.Errorf(
			"%w: password is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.AppID == "" {
		return domain.LoginResponse{}, fmt.Errorf(
			"%w: app_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.Login(ctx, req)
}

func (uc *SSOUseCase) Refresh(
	ctx context.Context,
	req domain.RefreshRequest,
) (domain.RefreshResponse, error) {
	if req.RefreshToken == "" {
		return domain.RefreshResponse{}, fmt.Errorf(
			"%w: refresh_token is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.Refresh(ctx, req)
}

func (uc *SSOUseCase) LogoutAll(
	ctx context.Context,
	req domain.LogoutAllRequest,
) (domain.LogoutAllResponse, error) {
	if req.UserID == "" {
		return domain.LogoutAllResponse{}, fmt.Errorf(
			"%w: user_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.LogoutAll(ctx, req)
}

func (uc *SSOUseCase) RegisterApp(
	ctx context.Context,
	req domain.RegisterAppRequest,
) (domain.RegisterAppResponse, error) {
	if req.Name == "" {
		return domain.RegisterAppResponse{}, fmt.Errorf(
			"%w: name is required",
			domain.ErrInvalidArgument,
		)
	}
	if req.Secret == "" {
		return domain.RegisterAppResponse{}, fmt.Errorf(
			"%w: secret is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.RegisterApp(ctx, req)
}

func (uc *SSOUseCase) DeactivateApp(
	ctx context.Context,
	req domain.DeactivateAppRequest,
) (domain.DeactivateAppResponse, error) {
	if req.AppID == "" {
		return domain.DeactivateAppResponse{}, fmt.Errorf(
			"%w: app_id is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.DeactivateApp(ctx, req)
}

func (uc *SSOUseCase) Logout(
	ctx context.Context,
	req domain.LogoutRequest,
) (domain.LogoutResponse, error) {
	if req.RefreshToken == "" {
		return domain.LogoutResponse{}, fmt.Errorf(
			"%w: token is required",
			domain.ErrInvalidArgument,
		)
	}

	return uc.sso.Logout(ctx, req)
}
