package grpc

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/sso/v1"
)

type SSOAdapter struct {
	client ssov1.AuthServiceClient
}

func NewSSOAdapter(client ssov1.AuthServiceClient) *SSOAdapter {
	return &SSOAdapter{client: client}
}

func (a *SSOAdapter) Register(
	ctx context.Context,
	req domain.RegisterRequest,
) (domain.RegisterResponse, error) {
	res, err := a.client.Register(ctx, &ssov1.RegisterRequest{
		Email:      req.Email,
		Password:   req.Password,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: req.MiddleName,
	})
	if err != nil {
		return domain.RegisterResponse{}, mapGRPCError(err)
	}

	return domain.RegisterResponse{UserID: res.UserId}, nil
}

func (a *SSOAdapter) Login(
	ctx context.Context,
	req domain.LoginRequest,
) (domain.LoginResponse, error) {
	res, err := a.client.Login(ctx, &ssov1.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
		AppId:    req.AppID,
	})
	if err != nil {
		return domain.LoginResponse{}, mapGRPCError(err)
	}

	return domain.LoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *SSOAdapter) Refresh(
	ctx context.Context,
	req domain.RefreshRequest,
) (domain.RefreshResponse, error) {
	res, err := a.client.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return domain.RefreshResponse{}, mapGRPCError(err)
	}

	return domain.RefreshResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}, nil
}

func (a *SSOAdapter) Logout(
	ctx context.Context,
	req domain.LogoutRequest,
) (domain.LogoutResponse, error) {
	_, err := a.client.Logout(ctx, &ssov1.LogoutRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return domain.LogoutResponse{}, mapGRPCError(err)
	}

	return domain.LogoutResponse{}, nil
}

func (a *SSOAdapter) LogoutAll(
	ctx context.Context,
	req domain.LogoutAllRequest,
) (domain.LogoutAllResponse, error) {
	_, err := a.client.LogoutAll(ctx, &ssov1.LogoutAllRequest{
		UserId: req.UserID,
	})
	if err != nil {
		return domain.LogoutAllResponse{}, mapGRPCError(err)
	}

	return domain.LogoutAllResponse{}, nil
}

func (a *SSOAdapter) RegisterApp(
	ctx context.Context,
	req domain.RegisterAppRequest,
) (domain.RegisterAppResponse, error) {
	res, err := a.client.RegisterApp(ctx, &ssov1.RegisterAppRequest{
		Name:        req.Name,
		Secret:      req.Secret,
		Description: req.Description,
	})
	if err != nil {
		return domain.RegisterAppResponse{}, mapGRPCError(err)
	}

	return domain.RegisterAppResponse{AppID: res.AppId}, nil
}

func (a *SSOAdapter) DeactivateApp(
	ctx context.Context,
	req domain.DeactivateAppRequest,
) (domain.DeactivateAppResponse, error) {
	_, err := a.client.DeactivateApp(ctx, &ssov1.DeactivateAppRequest{
		AppId: req.AppID,
	})
	if err != nil {
		return domain.DeactivateAppResponse{}, mapGRPCError(err)
	}

	return domain.DeactivateAppResponse{}, nil
}
