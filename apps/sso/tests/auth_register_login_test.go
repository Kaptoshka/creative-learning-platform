package sso_test

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/Kaptoshka/creative-learning-platform/sso-service/tests/suite"

	ssov1 "github.com/Kaptoshka/creative-learning-platform/libs/protos/gen/go/sso/v1"
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = ""
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	middleName := gofakeit.FirstName()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:      email,
		Password:   pass,
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
	})
	require.NoError(t, err)

	userID := respReg.GetUserId()
	assert.NotEmpty(t, userID)
	_, err = uuid.Parse(userID)
	require.NoError(t, err, "user_id must be a valid UUID")

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    st.AppID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	accessToken := respLogin.GetAccessToken()
	refreshToken := respLogin.GetRefreshToken()
	require.NotEmpty(t, accessToken)
	require.NotEmpty(t, refreshToken)

	publicKey := fetchJWKSPublicKey(t, st.JWKSUrl)

	tokenParsed, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodRSA)
		require.True(t, ok, "unexpected signing method")
		return publicKey, nil
	})
	require.NoError(t, err)
	require.True(t, tokenParsed.Valid)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, userID, claims["sub"].(string))
	assert.Equal(t, email, claims["email"].(string))

	assert.Empty(t, claims["app_id"])

	assert.Equal(t, "student", claims["role"].(string))

	assert.NotEmpty(t, claims["scope"])

	_, hasGroups := claims["group_ids"]
	assert.True(t, hasGroups)

	const deltaSeconds = 10
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"], deltaSeconds)
}

func TestRegisterLogin_Refresh_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:      email,
		Password:   pass,
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		MiddleName: gofakeit.FirstName(),
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    st.AppID,
	})
	require.NoError(t, err)

	oldRefresh := respLogin.GetRefreshToken()
	require.NotEmpty(t, oldRefresh)

	respRefresh, err := st.AuthClient.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: oldRefresh,
	})
	require.NoError(t, err)

	newAccessToken := respRefresh.GetAccessToken()
	newRefreshToken := respRefresh.GetRefreshToken()
	require.NotEmpty(t, newAccessToken)
	require.NotEmpty(t, newRefreshToken)

	assert.NotEqual(t, oldRefresh, newRefreshToken)

	_, err = st.AuthClient.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: oldRefresh,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "refresh token revoked")
}

func TestRegisterLogin_Logout_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:      email,
		Password:   pass,
		FirstName:  gofakeit.FirstName(),
		LastName:   gofakeit.LastName(),
		MiddleName: gofakeit.FirstName(),
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    st.AppID,
	})
	require.NoError(t, err)

	refreshToken := respLogin.GetRefreshToken()

	_, err = st.AuthClient.Logout(ctx, &ssov1.LogoutRequest{
		RefreshToken: refreshToken,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Refresh(ctx, &ssov1.RefreshRequest{
		RefreshToken: refreshToken,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "refresh token revoked")
}

func TestRegisterLogin_DuplicateRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	middleName := gofakeit.FirstName()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:      email,
		Password:   pass,
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:      email,
		Password:   pass,
		FirstName:  firstName,
		LastName:   lastName,
		MiddleName: middleName,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	require.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		firstName   string
		lastName    string
		middleName  string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			middleName:  gofakeit.FirstName(),
			expectedErr: "password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			middleName:  gofakeit.FirstName(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Empty First Name",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			firstName:   "",
			lastName:    gofakeit.LastName(),
			middleName:  gofakeit.FirstName(),
			expectedErr: "first_name is required",
		},
		{
			name:        "Register with Empty Last Name",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			firstName:   gofakeit.FirstName(),
			lastName:    "",
			middleName:  gofakeit.FirstName(),
			expectedErr: "last_name is required",
		},
		{
			name:        "Register with Empty Name",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			firstName:   "",
			lastName:    "",
			middleName:  "",
			expectedErr: "first_name is required",
		},
		{
			name:        "Register with Both Empty Email and Password",
			email:       "",
			password:    "",
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			middleName:  gofakeit.FirstName(),
			expectedErr: "email is required",
		},
		{
			name:        "Register with Full Empty",
			email:       "",
			password:    "",
			firstName:   "",
			lastName:    "",
			middleName:  "",
			expectedErr: "email is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:      tt.email,
				Password:   tt.password,
				FirstName:  tt.firstName,
				LastName:   tt.lastName,
				MiddleName: tt.middleName,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       string
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       st.AppID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			appID:       st.AppID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       st.AppID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       st.AppID,
			expectedErr: "invalid email or password",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
		{
			name:        "Login with Invalid AppID format",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       "not-a-uuid",
			expectedErr: "app_id must be a valid UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:      gofakeit.Email(),
				Password:   randomFakePassword(),
				FirstName:  gofakeit.FirstName(),
				LastName:   gofakeit.LastName(),
				MiddleName: gofakeit.FirstName(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func fetchJWKSPublicKey(t *testing.T, jwksURL string) *rsa.PublicKey {
	t.Helper()

	resp, err := http.Get(jwksURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			N   string `json:"n"`
			E   string `json:"e"`
			Kid string `json:"kid"`
			Alg string `json:"alg"`
		} `json:"keys"`
	}

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	require.NoError(t, err)
	require.NotEmpty(t, jwks.Keys, "JWKS must contain at least one key")

	key := jwks.Keys[0]

	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	require.NoError(t, err)

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	require.NoError(t, err)

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
