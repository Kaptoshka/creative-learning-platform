package domain

import (
	"errors"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
	ErrRoleNotFound = errors.New("role not found")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
)

var (
	ErrTokenNotFound = errors.New("refresh token not found")
	ErrTokenExpired  = errors.New("refresh token expired")
	ErrTokenRevoked  = errors.New("refresh token revoked")
)

var ErrAppExists = errors.New("app already exists")
