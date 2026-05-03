package driven

import (
	"context"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/google/uuid"
)

type UserRepo interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
		firstName string,
		lastName string,
		middleName string,
	) (uuid.UUID, error)
	User(
		ctx context.Context,
		email string,
	) (models.User, error)
	UserByID(
		ctx context.Context,
		userID uuid.UUID,
	) (models.User, error)
}

type RoleRepo interface {
	LinkUserRole(
		ctx context.Context,
		userID uuid.UUID,
		roleID uuid.UUID,
	) error
	UserRole(
		ctx context.Context,
		userID uuid.UUID,
	) (string, error)
	RoleID(
		ctx context.Context,
		role string,
	) (uuid.UUID, error)
	Scope(
		ctx context.Context,
		userID uuid.UUID,
	) ([]string, error)
}

type AppRepo interface {
	App(
		ctx context.Context,
		appID uuid.UUID,
	) (models.App, error)
	DeactivateApp(
		ctx context.Context,
		appID uuid.UUID,
	) error
	ActivateApp(
		ctx context.Context,
		appID uuid.UUID,
	) error
	RegisterApp(
		ctx context.Context,
		name string,
		secret string,
		description string,
	) (uuid.UUID, error)
}

type GroupRepo interface {
	UserGroups(
		ctx context.Context,
		userID uuid.UUID,
	) ([]uuid.UUID, error)
}

type RefreshRepo interface {
	Save(
		ctx context.Context,
		userID uuid.UUID,
		tokenHash string,
		expiresAt time.Time,
	) (uuid.UUID, error)

	ByHash(
		ctx context.Context,
		tokenHash string,
	) (models.RefreshToken, error)

	Revoke(
		ctx context.Context,
		tokenID uuid.UUID,
	) error

	RevokeAll(
		ctx context.Context,
		userID uuid.UUID,
	) error

	DeleteExpired(
		ctx context.Context,
	) error
}
