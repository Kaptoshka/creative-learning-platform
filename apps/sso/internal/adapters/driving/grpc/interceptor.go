package grpc

import (
	"context"
	"crypto/rsa"
	"strings"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AdminOnlyInterceptor(publicKey *rsa.PublicKey) grpc.UnaryServerInterceptor {
	adminMethods := map[string]bool{
		"/sso.v1.AuthService/RegisterApp":   true,
		"/sso.v1.AuthService/DeactivateApp": true,
		"/sso.v1.AuthService/ActivateApp":   true,
	}
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if !adminMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		token, err := extractToken(ctx)
		if err != nil {
			return nil, err
		}

		claims, err := parseToken(token, publicKey)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		role, ok := claims["role"].(string)
		if !ok || role != "admin" {
			return nil, status.Error(codes.PermissionDenied, "admin role required")
		}

		return handler(ctx, req)
	}
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "metadata is required")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "authorization header is required")
	}

	half := 2
	parts := strings.SplitN(values[0], " ", half)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", status.Error(codes.Unauthenticated, "authorization header must be Bearer <token>")
	}

	return parts[1], nil
}

func parseToken(tokenStr string, publicKey *rsa.PublicKey) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"unexpected signing method: %v", t.Header["alg"],
			)
		}
		return publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid claims")
	}

	return claims, nil
}
