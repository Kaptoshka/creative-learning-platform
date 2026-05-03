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

type Service struct {
	log         *slog.Logger
	appRepo     driven.AppRepo
	roleRepo    driven.RoleRepo
	userRepo    driven.UserRepo
	groupRepo   driven.GroupRepo
	refreshRepo driven.RefreshRepo
	tokens      TokenService
	refreshTTL  time.Duration
}

// New returns a new instance of Service.
func New(
	log *slog.Logger,
	appRepo driven.AppRepo,
	roleRepo driven.RoleRepo,
	userRepo driven.UserRepo,
	groupRepo driven.GroupRepo,
	refreshRepo driven.RefreshRepo,
	tokens TokenService,
	refreshTTL time.Duration,
) Service {
	return Service{
		log:         log,
		appRepo:     appRepo,
		roleRepo:    roleRepo,
		userRepo:    userRepo,
		groupRepo:   groupRepo,
		refreshRepo: refreshRepo,
		tokens:      tokens,
		refreshTTL:  refreshTTL,
	}
}

// Login service layer functions that
// checks if user with given credentials exists in the system
//
// If user exists, but password incorrect, returns ErrInvalidCredentials.
// If user does not exist, returns ErrUserNotFound.
func (a *Service) Login(
	ctx context.Context,
	email string,
	password string,
	appID uuid.UUID,
) (string, string, error) {
	const op = "services.auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)
	log.DebugContext(ctx, "attempting to login user")

	user, err := a.userRepo.User(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			log.WarnContext(ctx, "user not found", slog.Any("error", err))

			return "", "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
		}
		log.ErrorContext(ctx, "failed to get user", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.WarnContext(ctx, "invalid credentials", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, domain.ErrInvalidCredentials)
	}

	log.DebugContext(ctx, "credentials valid")

	_, err = a.appRepo.App(ctx, appID)
	if err != nil {
		if errors.Is(err, domain.ErrAppNotFound) {
			log.WarnContext(ctx, "app not found", slog.Any("error", err))
			return "", "", fmt.Errorf("%s: %w", op, domain.ErrInvalidAppID)
		}
		log.ErrorContext(ctx, "failed to get app", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.DebugContext(ctx, "app info found")

	role, err := a.roleRepo.UserRole(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user role", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	scope, err := a.roleRepo.Scope(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user permission scope", slog.Any("error", err))

		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	groupIDs, err := a.groupRepo.UserGroups(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user groups", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := a.tokens.GenerateNewToken(user, role, scope, groupIDs)
	if err != nil {
		log.ErrorContext(ctx, "failed to generate token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	rawRefresh, tokenHash := generateRefreshToken()

	_, err = a.refreshRepo.Save(
		ctx,
		user.ID,
		tokenHash,
		time.Now().Add(a.refreshTTL),
	)
	if err != nil {
		log.ErrorContext(ctx, "failed to save refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "user logged in", slog.String("user_id", user.ID.String()))

	return accessToken, rawRefresh, nil
}

// RegisterNewUser service layer function that implements user registration.
func (a *Service) RegisterNewUser(
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
	log.DebugContext(ctx, "registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.ErrorContext(ctx, "failed to generate password hash", slog.Any("error", err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userRepo.SaveUser(
		ctx, email, passHash, firstName, lastName, middleName,
	)
	if err != nil {
		if errors.Is(err, domain.ErrUserExists) {
			log.WarnContext(ctx, "failed to save user", slog.Any("error", domain.ErrUserExists))
			return uuid.Nil, fmt.Errorf("%s: %w", op, domain.ErrUserExists)
		}
		log.ErrorContext(ctx, "failed to save user", slog.Any("error", err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	roleID, err := a.roleRepo.RoleID(ctx, permissions.RoleStudent)
	if err != nil {
		log.ErrorContext(ctx, "failed to get role id", slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	err = a.roleRepo.LinkUserRole(ctx, id, roleID)
	if err != nil {
		log.ErrorContext(ctx, "failed to link user role", slog.Any("error", err))

		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "user registered")

	return id, nil
}

func (a *Service) RefreshToken(
	ctx context.Context,
	rawRefreshToken string,
) (string, string, error) {
	const op = "services.auth.RefreshToken"

	log := a.log.With(slog.String("op", op))
	log.DebugContext(ctx, "refreshing token")

	tokenHash := hashToken(rawRefreshToken)

	stored, err := a.refreshRepo.ByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, domain.ErrTokenNotFound) {
			return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenNotFound)
		}
		log.ErrorContext(ctx, "failed to get refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if stored.RevokedAt != nil {
		log.WarnContext(ctx, "refresh token already revoked",
			slog.String("token_id", stored.ID.String()),
		)
		_ = a.refreshRepo.RevokeAll(ctx, stored.UserID)
		return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenRevoked)
	}

	if time.Now().After(stored.ExpiresAt) {
		return "", "", fmt.Errorf("%s: %w", op, domain.ErrTokenExpired)
	}

	if err = a.refreshRepo.Revoke(ctx, stored.ID); err != nil {
		log.ErrorContext(ctx, "failed to revoke old refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	user, err := a.userRepo.UserByID(ctx, stored.UserID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	role, err := a.roleRepo.UserRole(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user role", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	scope, err := a.roleRepo.Scope(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user scope", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	groupIDs, err := a.groupRepo.UserGroups(ctx, user.ID)
	if err != nil {
		log.ErrorContext(ctx, "failed to get user groups", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	accessToken, err := a.tokens.GenerateNewToken(
		user, role, scope, groupIDs,
	)
	if err != nil {
		log.ErrorContext(ctx, "failed to generate access token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	rawRefresh, newHash := generateRefreshToken()

	_, err = a.refreshRepo.Save(ctx, user.ID, newHash, time.Now().Add(a.refreshTTL))
	if err != nil {
		log.ErrorContext(ctx, "failed to save new refresh token", slog.Any("error", err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "token refreshed", slog.String("user_id", user.ID.String()))

	return accessToken, rawRefresh, nil
}

func (a *Service) Logout(
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
		log.ErrorContext(ctx, "failed to get refresh token", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = a.refreshRepo.Revoke(ctx, stored.ID); err != nil {
		log.ErrorContext(ctx, "failed to revoke refresh token", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "user logged out", slog.String("user_id", stored.UserID.String()))

	return nil
}

func (a *Service) LogoutAll(
	ctx context.Context,
	userID uuid.UUID,
) error {
	const op = "services.auth.LogoutAll"

	log := a.log.With(slog.String("op", op))

	if err := a.refreshRepo.RevokeAll(ctx, userID); err != nil {
		log.ErrorContext(ctx, "failed to revoke all refresh tokens", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.InfoContext(ctx, "all sessions revoked", slog.String("user_id", userID.String()))

	return nil
}

func generateRefreshToken() (string, string) {
	byte32 := 32
	b := make([]byte, byte32)
	_, _ = rand.Read(b)
	raw := hex.EncodeToString(b)
	hash := hashToken(raw)
	return raw, hash
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func (a *Service) RegisterApp(
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

	a.log.InfoContext(ctx, "app registered",
		slog.String("app_id", id.String()),
		slog.String("name", name),
	)

	return id, nil
}

func (a *Service) DeactivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "services.auth.DeactivateApp"

	if err := a.appRepo.DeactivateApp(ctx, appID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.InfoContext(ctx, "app deactivated", slog.String("app_id", appID.String()))

	return nil
}

func (a *Service) ActivateApp(
	ctx context.Context,
	appID uuid.UUID,
) error {
	const op = "services.auth.ActivateApp"

	if err := a.appRepo.ActivateApp(ctx, appID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.InfoContext(ctx, "app activated", slog.String("app_id", appID.String()))

	return nil
}
