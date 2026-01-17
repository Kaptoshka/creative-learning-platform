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

func GetUserID(ctx context.Context) (int64, error) {
	val, ok := ctx.Value(contextKeyUserID).(int64)
	if !ok {
		return 0, errors.New("user id not found in context")
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

func WithUser(ctx context.Context, userID int64, role string) context.Context {
	ctx = context.WithValue(ctx, contextKeyUserID, userID)
	ctx = context.WithValue(ctx, contextKeyRole, role)
	return ctx
}
