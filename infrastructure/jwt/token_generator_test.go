package jwt_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestJwtAccessTokenGenerator(t *testing.T) {
	testCases := []struct {
		description       string
		argUserID         int64
		expectedSub       string
		expectedTokenType string
	}{
		{
			description:       "ShouldGenerateAccessToken",
			argUserID:         10,
			expectedSub:       "10",
			expectedTokenType: "access",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cfg := config.DefaultConfig()
			secretKeyBytes := []byte(cfg.Jwt.SecretKey)

			tokenGeneratorFunc := jwt.NewJwtAccessTokenGeneratorFunc(cfg)
			tok, err := tokenGeneratorFunc(tc.argUserID)

			assert.NoError(t, err)

			var actualClaims jwtGo.MapClaims
			jwtGo.ParseWithClaims(
				tok,
				&actualClaims,
				func(token *jwtGo.Token) (interface{}, error) {
					return secretKeyBytes, nil
				},
			)

			assert.Equal(t, tc.expectedSub, actualClaims["sub"])
			assert.Equal(t, tc.expectedTokenType, actualClaims["type"])
		})
	}
}

func TestJwtRefreshTokenGenerator(t *testing.T) {
	testCases := []struct {
		description       string
		argUserID         int64
		expectedSub       string
		expectedTokenType string
	}{
		{
			description:       "ShouldGenerateRefreshToken",
			argUserID:         10,
			expectedSub:       "10",
			expectedTokenType: "refresh",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			cfg := config.DefaultConfig()
			secretKeyBytes := []byte(cfg.Jwt.SecretKey)

			tokenGeneratorFunc := jwt.NewJwtRefreshTokenGeneratorFunc(cfg)
			tok, err := tokenGeneratorFunc(tc.argUserID)

			assert.NoError(t, err)

			var actualClaims jwtGo.MapClaims
			jwtGo.ParseWithClaims(
				tok,
				&actualClaims,
				func(token *jwtGo.Token) (interface{}, error) {
					return secretKeyBytes, nil
				},
			)

			assert.Equal(t, tc.expectedSub, actualClaims["sub"])
			assert.Equal(t, tc.expectedTokenType, actualClaims["type"])
		})
	}
}
