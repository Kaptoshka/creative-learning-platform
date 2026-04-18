package auth

import (
	"context"
	"errors"
)

type contextKey string

const (
	contextKeyUserID contextKey = "user_id"
	contextKeyRole   contextKey = "role"
)

const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
	RoleDev     = "dev"
)

func GetUserID(ctx context.Context) (string, error) {
	val, ok := ctx.Value(contextKeyUserID).(string)
	if !ok || val == "" {
		return "", errors.New("user id not found in context")
	}
	return val, nil
}

func GetUserRole(ctx context.Context) string {
	val, ok := ctx.Value(contextKeyRole).(string)
	if !ok {
		return ""
	}
	return val
}

func New(
	ctx context.Context,
	userID string,
	role string,
) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID, userID)
	ctx = context.WithValue(ctx, contextKeyRole, role)
	return ctx
}
