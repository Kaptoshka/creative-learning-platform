package grpc

import (
	"context"
	"errors"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/ports/driving"

	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/proto/sso/v1"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyValue = 0
)

type serverAPI struct {
	ssov1.UnimplementedAuthServiceServer
	auth driving.AuthService
}

func Register(gRPC *grpc.Server, auth driving.AuthService) {
	ssov1.RegisterAuthServiceServer(gRPC, &serverAPI{auth: auth})
}

// Login implements login of the user in SSO
func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	appID, err := uuid.Parse(req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "app_id must be a valid UUID")
	}

	accessToken, refreshToken, err := s.auth.Login(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
		appID,
	)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}
		if errors.Is(err, domain.ErrInvalidAppID) {
			return nil, status.Error(codes.NotFound, "app not found")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Register implements registration of the user in SSO
func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
		req.GetFirstName(),
		req.GetLastName(),
		req.GetMiddleName(),
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID.String(),
	}, nil
}

func (s *serverAPI) Refresh(
	ctx context.Context,
	req *ssov1.RefreshRequest,
) (*ssov1.RefreshResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	accessToken, refreshToken, err := s.auth.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, domain.ErrTokenExpired) {
			return nil, status.Error(codes.Unauthenticated, "refresh token expired")
		}
		if errors.Is(err, domain.ErrTokenRevoked) {
			return nil, status.Error(codes.Unauthenticated, "refresh token revoked")
		}
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil, status.Error(codes.Unauthenticated, "refresh token not found")
		}

		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &ssov1.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) Logout(
	ctx context.Context,
	req *ssov1.LogoutRequest,
) (*ssov1.LogoutResponse, error) {
	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	if err := s.auth.Logout(ctx, req.GetRefreshToken()); err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return nil, status.Error(codes.NotFound, "refresh token not found")
		}

		return nil, status.Error(codes.Internal, "failed to logout")
	}

	return &ssov1.LogoutResponse{}, nil
}

func (s *serverAPI) LogoutAll(
	ctx context.Context,
	req *ssov1.LogoutAllRequest,
) (*ssov1.LogoutAllResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "user_id must be a valid UUID")
	}

	if err := s.auth.LogoutAll(ctx, userID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to logout all sessions")
	}

	return &ssov1.LogoutAllResponse{}, nil
}

func (s *serverAPI) RegisterApp(
	ctx context.Context,
	req *ssov1.RegisterAppRequest,
) (*ssov1.RegisterAppResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetSecret() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret is required")
	}

	appID, err := s.auth.RegisterApp(ctx, req.GetName(), req.GetSecret(), req.GetDescription())
	if err != nil {
		if errors.Is(err, domain.ErrAppExists) {
			return nil, status.Error(codes.AlreadyExists, "app already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register app")
	}

	return &ssov1.RegisterAppResponse{AppId: appID.String()}, nil
}

func (s *serverAPI) DeactivateApp(
	ctx context.Context,
	req *ssov1.DeactivateAppRequest,
) (*ssov1.DeactivateAppResponse, error) {
	appID, err := uuid.Parse(req.GetAppId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "app_id must be a valid UUID")
	}

	if err := s.auth.DeactivateApp(ctx, appID); err != nil {
		return nil, status.Error(codes.Internal, "failed to deactivate app")
	}

	return &ssov1.DeactivateAppResponse{}, nil
}

// validateLogin validates the login request
// Email, password and AppId must be provided.
// If not it returns an error.
func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == "" {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

// validateRegister validates the register request
// Email, password, first_name, last_name must be provided.
// If not it returns an error.
func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetFirstName() == "" {
		return status.Error(codes.InvalidArgument, "first_name is required")
	}

	if req.GetLastName() == "" {
		return status.Error(codes.InvalidArgument, "last_name is required")
	}

	return nil
}
