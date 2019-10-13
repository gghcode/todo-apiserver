package auth_test

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/app/api/auth"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/stretchr/testify/suite"
)

type MiddlewareAdapterUnit struct {
	suite.Suite

	cfg      config.Configuration
	jwtParam auth.JwtParam

	fakeTokenRepo fake.TokenRepository
	fakeUserRepo  fake.UserRepository
	fakePassport  fake.Passport
	service       auth.Service
}

func TestAuthMiddlewareAdapterUnit(t *testing.T) {
	suite.Run(t, new(MiddlewareAdapterUnit))
}

func (suite *ServiceUnit) TestVerifyAccessToken() {
	var fakeUserID int64 = 10

	fakeAccessToken, _ := auth.CreateAccessToken(suite.jwtParam, fakeUserID)
	fakeJwtClaims, _ := auth.ExtractTokenClaims(suite.jwtParam, fakeAccessToken)

	testCases := []struct {
		description    string
		argSecret      string
		argAccessToken string
		expected       jwt.MapClaims
		expectedErr    error
	}{
		{
			description:    "ShouldVerificationSuccess",
			argSecret:      suite.cfg.Jwt.SecretKey,
			argAccessToken: fmt.Sprintf("Bearer %s", fakeAccessToken),
			expected:       fakeJwtClaims,
			expectedErr:    nil,
		},
		{
			description:    "ShouldBeUnauthorizedWhenNotContainToken",
			argSecret:      suite.cfg.Jwt.SecretKey,
			argAccessToken: "",
			expected:       nil,
			expectedErr:    auth.ErrNotContainToken,
		},
		{
			description:    "ShouldBeUnauthorizedWhenWrongTokenInfo",
			argSecret:      suite.cfg.Jwt.SecretKey,
			argAccessToken: "accessToken",
			expected:       nil,
			expectedErr:    auth.ErrInvalidToken,
		},
		{
			description:    "ShouldBeUnauthorizedWhenInvalidTokenType",
			argSecret:      suite.cfg.Jwt.SecretKey,
			argAccessToken: "abcd accessToken",
			expected:       nil,
			expectedErr:    auth.ErrInvalidTokenType,
		},
		{
			description:    "ShouldBeUnauthorizedWhenInvalidAccessToken",
			argSecret:      suite.cfg.Jwt.SecretKey,
			argAccessToken: "Bearer fasdfasdgvad",
			expected:       nil,
			expectedErr:    auth.ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := auth.VerifyAccessToken(tc.argSecret, tc.argAccessToken)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
