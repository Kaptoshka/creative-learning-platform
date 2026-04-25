package driven

import (
	"context"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"
)

type UserRepo interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
		firstName string,
		lastName string,
		middleName string,
	) (int64, error)
	User(
		ctx context.Context,
		email string,
	) (models.User, error)
}

type RoleRepo interface {
	LinkUserRole(
		ctx context.Context,
		userID int64,
		roleID int64,
	) error
	UserRole(
		ctx context.Context,
		userID int64,
	) (string, error)
	RoleID(
		ctx context.Context,
		role string,
	) (int64, error)
	Scope(
		ctx context.Context,
		userID int64,
	) ([]string, error)
}

type AppRepo interface {
	App(
		ctx context.Context,
		appID int,
	) (models.App, error)
}
