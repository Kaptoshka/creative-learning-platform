package storage

import (
	"context"
	"errors"

	"sso/internal/domain/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
	ErrRoleNotFound = errors.New("role not found")
)

type UserStorage interface {
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

type RoleStorage interface {
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

type AppStorage interface {
	App(
		ctx context.Context,
		appID int,
	) (models.App, error)
}
