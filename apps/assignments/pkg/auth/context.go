package auth

import (
	"context"
	"errors"
)

type contextKey string

const (
	ContextKeyUserID  contextKey = "user_id"
	ContextKeyRole    contextKey = "role"
	ContextKeyGroupID contextKey = "group_id"
)

const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
	RoleDev     = "dev"
)

func GetUserID(ctx context.Context) (string, error) {
	val, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok || val == "" {
		return "", errors.New("user id not found in context")
	}
	return val, nil
}

func GetUserRole(ctx context.Context) string {
	val, ok := ctx.Value(ContextKeyRole).(string)
	if !ok {
		return ""
	}
	return val
}

func GetGroupID(ctx context.Context) string {
	val, ok := ctx.Value(ContextKeyGroupID).(string)
	if !ok {
		return ""
	}
	return val
}
