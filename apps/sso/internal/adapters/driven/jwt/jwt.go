package jwt

import (
	"strings"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	ttl time.Duration
}

func New(ttl time.Duration) *JWTService {
	return &JWTService{
		ttl: ttl,
	}
}

// GenerateNewToken generates a new JWT token
// for the given user, app, duration, role, and permission scope.
func (j *JWTService) GenerateNewToken(
	user models.User,
	app models.App,
	role string,
	scope []string,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(j.ttl).Unix()
	claims["role"] = role
	claims["scope"] = strings.Join(scope, " ")

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
