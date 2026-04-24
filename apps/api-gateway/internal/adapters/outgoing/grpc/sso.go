package grpc

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/gateway/internal/core/domain"
	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/proto/sso/v1"
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

	return domain.LoginResponse{Token: res.Token}, nil
}

func (a *SSOAdapter) Logout(
	ctx context.Context,
	req domain.LogoutRequest,
) (domain.LogoutResponse, error) {
	res, err := a.client.Logout(ctx, &ssov1.LogoutRequest{
		Token: req.Token,
	})
	if err != nil {
		return domain.LogoutResponse{}, mapGRPCError(err)
	}

	return domain.LogoutResponse{Success: res.Success}, nil
}
