package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kaptoshka/creative-learning-platform/assignment-service/pkg/auth"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtSecret string
}

func NewAuthInterceptor(jwtSecret string) *AuthInterceptor {
	return &AuthInterceptor{jwtSecret: jwtSecret}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		claims, err := a.parseToken(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx = context.WithValue(ctx, auth.ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, auth.ContextKeyRole, claims.Role)
		ctx = context.WithValue(ctx, auth.ContextKeyGroupID, claims.GroupID)

		return handler(ctx, req)
	}
}

type Claims struct {
	UserID  uuid.UUID
	Role    string
	GroupID uuid.UUID
}

func (a *AuthInterceptor) parseToken(ctx context.Context) (*Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("invalid authorization format")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	claims, err := mapToClaims(mapClaims)
	if err != nil {
		return nil, fmt.Errorf("map claims: %w", err)
	}

	return claims, nil
}

func mapToClaims(mc jwt.MapClaims) (*Claims, error) {
	userIDStr, ok := mc["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing user_id claim")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	role, _ := mc["role"].(string)

	var groupID uuid.UUID
	if groupIDStr, ok := mc["group_id"].(string); ok && groupIDStr != "" {
		groupID, err = uuid.Parse(groupIDStr)
		if err != nil {
			return nil, fmt.Errorf("invalid group_id: %w", err)
		}
	}

	return &Claims{
		UserID:  userID,
		Role:    role,
		GroupID: groupID,
	}, nil
}
