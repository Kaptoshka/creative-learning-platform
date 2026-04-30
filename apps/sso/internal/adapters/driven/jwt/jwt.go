package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/internal/core/domain/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	ttl        time.Duration
	privateKey *rsa.PrivateKey
	keyID      string
}

func New(
	ttl time.Duration,
	privateKey *rsa.PrivateKey,
	keyID string,
) *JWTService {
	return &JWTService{
		ttl:        ttl,
		privateKey: privateKey,
		keyID:      keyID,
	}
}

// GenerateNewToken generates a new JWT token
// for the given user, app, duration, role, and permission scope.
func (j *JWTService) GenerateNewToken(
	user models.User,
	role string,
	scope []string,
	groupIDs []uuid.UUID,
) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	token.Header["kid"] = j.keyID

	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()

	claims["sub"] = user.ID.String()
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(j.ttl).Unix()

	claims["email"] = user.Email
	claims["role"] = role
	claims["scope"] = strings.Join(scope, " ")

	groups := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		groups[i] = id.String()
	}
	claims["group_ids"] = groups

	return token.SignedString(j.privateKey)
}

func (j *JWTService) PublicKey() *rsa.PublicKey {
	return &j.privateKey.PublicKey
}

func (j *JWTService) KeyID() string {
	return j.keyID
}

func ParseRSAPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("PKCS8 key is not RSA")
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}
}
