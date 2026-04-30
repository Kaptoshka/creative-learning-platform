package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/permissions"
	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/ports/driven"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type TokenService interface {
	GenerateNewToken(
		user models.User,
		role string,
		scope []string,
		groupIDs []uuid.UUID,
	) (string, error)
}

type authService struct {
	log         *slog.Logger
	appRepo     driven.AppRepo
	roleRepo    driven.RoleRepo
	userRepo    driven.UserRepo
	groupRepo   driven.GroupRepo
	refreshRepo driven.RefreshRepo
	tokens      TokenService
	refreshTTL  time.Duration
}

// New returns a new instance of authService
func New(
	log *slog.Logger,
	appRepo driven.AppRepo,
	roleRepo driven.RoleRepo,
	userRepo driven.UserRepo,
	groupRepo driven.GroupRepo,
	refreshRepo driven.RefreshRepo,
	tokens TokenService,
	refreshTTL time.Duration,
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
	appID uuid.UUID,
) (string, string, error) {
	const op = "services.auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Debug("attempting to login user")

	user, err := a.userRepo.User(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.Warn("user not found", slog.Any("error", err))

			return "", "", fmt.Errorf("%s: %v", op, domain.ErrInvalidCredentials)
		}
		log.Error("failed to get user", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %v", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %v", op, domain.ErrInvalidCredentials)
	}

	log.Debug("credentials valid")

	_, err = a.appRepo.App(ctx, appID)
	if err != nil {
		if errors.Is(err, domain.ErrAppNotFound) {
			log.Warn("app not found", slog.Any("error", err))
			return "", "", fmt.Errorf("%s: %v", op, domain.ErrInvalidAppID)
		}
		log.Error("failed to get app", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %v", op, err)
	}

	log.Debug("app info found")

	role, err := a.roleRepo.UserRole(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user role", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %v", op, err)
	}

	scope, err := a.roleRepo.Scope(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user permission scope", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %v", op, err)
	}

	groupIDs, err := a.groupRepo.UserGroups(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user groups", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := a.tokens.GenerateNewToken(user, role, scope, groupIDs)
	if err != nil {
		log.Error("failed to generate token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %v", op, err)
	}

	rawRefresh, tokenHash := generateRefreshToken()

	_, err = a.refreshRepo.Save(
		ctx,
		user.ID,
		tokenHash,
		time.Now().Add(a.refreshTTL),
	)
	if err != nil {
		log.Error("failed to save refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in", slog.String("user_id", user.ID.String()))

	return accessToken, rawRefresh, nil
}

// RegisterNewUser service layer function that implements user registration
func (a *authService) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
	firstName string,
	lastName string,
	middleName string,
) (uuid.UUID, error) {
	const op = "services.auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Debug("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("error", err))
		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	id, err := a.userRepo.SaveUser(
		ctx, email, passHash, firstName, lastName, middleName,
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			log.Warn("failed to save user", slog.Any("error", domain.ErrUserExists))
			return uuid.Nil, fmt.Errorf("%s: %v", op, domain.ErrUserExists)
		}
		log.Error("failed to save user", slog.Any("error", err))
		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	roleID, err := a.roleRepo.RoleID(ctx, permissions.RoleStudent)
	if err != nil {
		log.Error("failed to get role id", slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	err = a.roleRepo.LinkUserRole(ctx, id, roleID)
	if err != nil {
		log.Error("failed to link user role", slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Info("user registered")

	return id, nil
}

func (a *authService) RefreshToken(
	ctx context.Context,
	rawRefreshToken string,
) (string, string, error) {
	const op = "services.auth.RefreshToken"

	log := a.log.With(slog.String("op", op))
	log.Debug("refreshing token")

	tokenHash := hashToken(rawRefreshToken)

	stored, err := a.refreshRepo.ByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		log.Error("failed to get refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if stored.RevokedAt != nil {
		log.Warn("refresh token already revoked",
			slog.String("token_id", stored.ID.String()),
		)
		_ = a.refreshRepo.RevokeAll(ctx, stored.UserID)
		return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenRevoked)
	}

	if time.Now().After(stored.ExpiresAt) {
		return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenExpired)
	}

	if err := a.refreshRepo.Revoke(ctx, stored.ID); err != nil {
		log.Error("failed to revoke old refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userRepo.UserByID(ctx, stored.UserID)
	if err != nil {
		log.Error("failed to get user", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	role, err := a.roleRepo.UserRole(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user role", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	scope, err := a.roleRepo.Scope(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user scope", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	groupIDs, err := a.groupRepo.UserGroups(ctx, user.ID)
	if err != nil {
		log.Error("failed to get user groups", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := a.tokens.GenerateNewToken(
		user, role, scope, groupIDs,
	)
	if err != nil {
		log.Error("failed to generate access token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	rawRefresh, newHash := generateRefreshToken()

	_, err = a.refreshRepo.Save(ctx, user.ID, newHash, time.Now().Add(a.refreshTTL))
	if err != nil {
		log.Error("failed to save new refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token refreshed", slog.String("user_id", user.ID.String()))

	return accessToken, rawRefresh, nil
}

func (a *authService) Logout(
	ctx context.Context,
	rawRefreshToken string,
) error {
	const op = "services.auth.Logout"

	log := a.log.With(slog.String("op", op))

	tokenHash := hashToken(rawRefreshToken)

	stored, err := a.refreshRepo.ByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		log.Error("failed to get refresh token", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := a.refreshRepo.Revoke(ctx, stored.ID); err != nil {
		log.Error("failed to revoke refresh token", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged out", slog.String("user_id", stored.UserID.String()))

	return nil
}

func (a *authService) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	const op = "services.auth.LogoutAll"

	log := a.log.With(slog.String("op", op))

	if err := a.refreshRepo.RevokeAll(ctx, userID); err != nil {
		log.Error("failed to revoke all refresh tokens", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("all sessions revoked", slog.String("user_id", userID.String()))

	return nil
}

func generateRefreshToken() (raw string, hash string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	raw = hex.EncodeToString(b)
	hash = hashToken(raw)
	return raw, hash
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (a *authService) RegisterApp(
	ctx context.Context,
	name string,
	secret string,
	description string,
) (uuid.UUID, error) {
	const op = "services.auth.RegisterApp"

	id, err := a.appRepo.RegisterApp(ctx, name, secret, description)
	if err != nil {
		if errors.Is(err, domain.ErrAppExists) {
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrAppExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("app registered",
		slog.String("app_id", id.String()),
		slog.String("name", name),
	)

	return id, nil
}

func (a *authService) DeactivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "services.auth.DeactivateApp"

	if err := a.appRepo.DeactivateApp(ctx, appID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("app deactivated", slog.String("app_id", appID.String()))

	return nil
}

func (a *authService) ActivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "services.auth.ActivateApp"

	if err := a.appRepo.ActivateApp(ctx, appID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("app activated", slog.String("app_id", appID.String()))

	return nil
}
