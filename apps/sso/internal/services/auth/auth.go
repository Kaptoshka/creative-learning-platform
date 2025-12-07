package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/domain/models"
	"sso/internal/domain/permissions"
	"sso/internal/lib/jwt"
	"sso/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHash []byte,
		firstName string,
		lastName string,
		middleName string,
	) (int64, error)
	LinkUserRole(
		ctx context.Context,
		userID int64,
		roleID int64,
	) error
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	UserRole(ctx context.Context, userID int64) (string, error)
	RoleID(ctx context.Context, role string) (int64, error)
	Scope(ctx context.Context, userID int64) ([]string, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = storage.ErrUserExists
)

// New returns a new instance of Auth service
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system
//
// If user exists, but password incorrect, returns ErrInvalidCredentials.
// If user does not exist, returns ErrUserNotFound.
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	const op = "services.auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("attempting to login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.Any("error", err))

			return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user found")

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
	}

	log.Debug("credentials valid")

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("app not found", slog.Any("error", err))

			return "", fmt.Errorf("%s: %v", op, ErrInvalidAppID)
		}

		a.log.Error("failed to get app", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("app info found")

	role, err := a.userProvider.UserRole(ctx, int64(user.ID))
	if err != nil {
		a.log.Error("failed to get user role", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user role found")

	scope, err := a.userProvider.Scope(ctx, int64(user.ID))
	if err != nil {
		a.log.Error("failed to get user permission scope", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user permission scope found")

	token, err := jwt.GenerateNewToken(user, app, a.tokenTTL, role, scope)
	if err != nil {
		a.log.Error("failed to generate token", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("jwt token generated")

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
	firstName string,
	lastName string,
	middleName string,
) (int64, error) {
	const op = "services.auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Debug("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("password hash generated")

	id, err := a.userSaver.SaveUser(
		ctx, email, passHash, firstName, lastName, middleName,
	)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("failed to save user", slog.Any("error", storage.ErrUserExists))

			return 0, fmt.Errorf("%s: %v", op, storage.ErrUserExists)
		}

		log.Error("failed to save user", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user saved")

	roleID, err := a.userProvider.RoleID(ctx, permissions.RoleStudent)
	if err != nil {
		log.Error("failed to get role id", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("role id found")

	err = a.userSaver.LinkUserRole(ctx, id, roleID)
	if err != nil {
		log.Error("failed to link user role", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user role linked")

	log.Info("user registered")

	return id, nil
}
