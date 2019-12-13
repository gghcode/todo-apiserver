package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"github.com/stretchr/testify/assert"
)

func TestJwtVerifyAccessToken(t *testing.T) {
	testCfg := config.Configuration{
		Jwt: config.JwtConfig{
			SecretKey: "testkey",
		},
	}

	testSecretKeyBytes := []byte(testCfg.Jwt.SecretKey)

	testCases := []struct {
		description    string
		argCfg         config.Configuration
		argAccessToken string
		expectedClaims middleware.TokenClaims
		expectedErr    error
	}{
		{
			description: "ShouldReturnOK",
			argCfg:      testCfg,
			argAccessToken: fmt.Sprintf("Bearer %s",
				jwtToken(t, testSecretKeyBytes, time.Now().Add(1*time.Hour)),
			),
			expectedClaims: middleware.TokenClaims{
				UserID: 5,
			},
			expectedErr: nil,
		},
		{
			description:    "ShouldReturnErrNotContainTokenInHeader",
			argAccessToken: "",
			expectedClaims: middleware.TokenClaims{},
			expectedErr:    auth.ErrNotContainTokenInHeader,
		},
		{
			description:    "ShouldReturnErrInvalidToken",
			argAccessToken: "dfadfasdfasdfasdfasdfsdfsdf",
			expectedClaims: middleware.TokenClaims{},
			expectedErr:    auth.ErrInvalidToken,
		},
		{
			description:    "ShouldReturnErrInvalidTokenType",
			argAccessToken: "JWT fasdfasdfasdfasdfasdfsdfasdf",
			expectedClaims: middleware.TokenClaims{},
			expectedErr:    auth.ErrInvalidTokenType,
		},
		{
			description:    "ShouldReturnErrInvalidToken",
			argAccessToken: "Bearer fasdfasdfasdfasdfasdfjklasdf",
			expectedClaims: middleware.TokenClaims{},
			expectedErr:    auth.ErrInvalidToken,
		},
		{
			description: "ShouldErrTokenExpired",
			argCfg:      testCfg,
			argAccessToken: fmt.Sprintf("Bearer %s",
				jwtToken(t, testSecretKeyBytes, time.Unix(10, 0)),
			),
			expectedClaims: middleware.TokenClaims{},
			expectedErr:    auth.ErrTokenExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			accessTokenHandlerFactory := NewJwtAccessTokenVerifyHandlerFactory(tc.argCfg)
			accessTokenHandler := accessTokenHandlerFactory.Create()

			actualClaims, actualErr := accessTokenHandler(tc.argAccessToken)

			assert.Equal(t, tc.expectedClaims, actualClaims)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func jwtToken(t *testing.T, secretKeyBytes []byte, expiresAt time.Time) string {
	claims := &jwt.StandardClaims{
		Subject:   "5",
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  time.Unix(1000, 0000).Unix(),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString(secretKeyBytes)

	assert.NoError(t, err)

	return tokenString
}
