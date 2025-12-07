package jwt

import (
	"strings"
	"time"

	"sso/internal/domain/models"

	"github.com/golang-jwt/jwt"
)

// GenerateNewToken generates a new JWT token
// for the given user, app, duration, role, and permission scope.
func GenerateNewToken(
	user models.User,
	app models.App,
	duration time.Duration,
	role string,
	scope []string,
) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["email"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["role"] = role
	claims["scope"] = strings.Join(scope, " ")

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
