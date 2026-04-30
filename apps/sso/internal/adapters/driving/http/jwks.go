package http

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
)

type JWKSProvider interface {
	PublicKey() *rsa.PublicKey
	KeyID() string
}

type jwk struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func newJWKSHandler(provider JWKSProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pub := provider.PublicKey()

		key := jwk{
			Kty: "RSA",
			Use: "sig",
			Alg: "RS256",
			Kid: provider.KeyID(),
			N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			E: base64.RawURLEncoding.EncodeToString(
				big.NewInt(int64(pub.E)).Bytes(),
			),
		}

		resp := jwks{Keys: []jwk{key}}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode jwks", http.StatusInternalServerError)
			return
		}
	}
}
