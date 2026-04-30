package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/permissions"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/ports/driven"

	"golang.org/x/crypto/bcrypt"
)

type TokenService interface {
	GenerateNewToken(
		user models.User,
		app models.App,
		role string,
		scope []string,
	) (string, error)
}

type authService struct {
	log      *slog.Logger
	appRepo  driven.AppRepo
	roleRepo driven.RoleRepo
	userRepo driven.UserRepo
	tokens   TokenService
}

// New returns a new instance of authService
func New(
	log *slog.Logger,
	appRepo driven.AppRepo,
	roleRepo driven.RoleRepo,
	userRepo driven.UserRepo,
	tokens TokenService,
) authService {
	return authService{
		log:      log,
		appRepo:  appRepo,
		roleRepo: roleRepo,
		userRepo: userRepo,
		tokens:   tokens,
	}
}

// Login service layer functions that
// checks if user with given credentials exists in the system
//
// If user exists, but password incorrect, returns ErrInvalidCredentials.
// If user does not exist, returns ErrUserNotFound.
func (a *authService) Login(
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

	user, err := a.userRepo.User(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			a.log.Warn("user not found", slog.Any("error", err))

			return "", fmt.Errorf("%s: %v", op, domain.ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user found")

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("invalid credentials", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, domain.ErrInvalidCredentials)
	}

	log.Debug("credentials valid")

	app, err := a.appRepo.App(ctx, appID)
	if err != nil {
		if errors.Is(err, domain.ErrAppNotFound) {
			a.log.Warn("app not found", slog.Any("error", err))

			return "", fmt.Errorf("%s: %v", op, domain.ErrInvalidAppID)
		}

		a.log.Error("failed to get app", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("app info found")

	role, err := a.roleRepo.UserRole(ctx, int64(user.ID))
	if err != nil {
		a.log.Error("failed to get user role", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user role found")

	scope, err := a.roleRepo.Scope(ctx, int64(user.ID))
	if err != nil {
		a.log.Error("failed to get user permission scope", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user permission scope found")

	token, err := a.tokens.GenerateNewToken(user, app, role, scope)
	if err != nil {
		a.log.Error("failed to generate token", slog.Any("error", err))

		return "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("jwt token generated")

	return token, nil
}

// RegisterNewUser service layer function that implements user registration
func (a *authService) RegisterNewUser(
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

	id, err := a.userRepo.SaveUser(
		ctx, email, passHash, firstName, lastName, middleName,
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			log.Warn("failed to save user", slog.Any("error", domain.ErrUserExists))

			return 0, fmt.Errorf("%s: %v", op, domain.ErrUserExists)
		}

		log.Error("failed to save user", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user saved")

	roleID, err := a.roleRepo.RoleID(ctx, permissions.RoleStudent)
	if err != nil {
		log.Error("failed to get role id", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("role id found")

	err = a.roleRepo.LinkUserRole(ctx, id, roleID)
	if err != nil {
		log.Error("failed to link user role", slog.Any("error", err))

		return 0, fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("user role linked")

	log.Info("user registered")

	return id, nil
}
