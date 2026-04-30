package driving

import (
	"context"

	"github.com/google/uuid"
)

type AuthService interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID uuid.UUID,
	) (accessToken string, refreshToken string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
		firstName string,
		lastName string,
		middleName string,
	) (userID uuid.UUID, err error)

	Logout(
		ctx context.Context,
		rawRefreshToken string,
	) error

	LogoutAll(
		ctx context.Context,
		userID uuid.UUID,
	) error

	RefreshToken(
		ctx context.Context,
		rawRefreshToken string,
	) (accessToken string, refreshToken string, err error)

	RegisterApp(
		ctx context.Context,
		name string,
		secret string,
		description string,
	) (appID uuid.UUID, err error)

	DeactivateApp(ctx context.Context, appID uuid.UUID) error
	ActivateApp(ctx context.Context, appID uuid.UUID) error
}
