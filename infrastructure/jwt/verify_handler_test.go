package jwt

import (
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/stretchr/testify/assert"
)

func TestJwtVerifyAccessToken(t *testing.T) {
	testSecretKeyBytes := []byte("testkey")

	testCases := []struct {
		description       string
		argSecretKeyBytes []byte
		argAccessToken    string
		expectedClaims    func(claims jwt.MapClaims) jwt.MapClaims
		expectedErr       error
	}{
		{
			description:       "ShouldReturnOK",
			argSecretKeyBytes: testSecretKeyBytes,
			argAccessToken: fmt.Sprintf("Bearer %s",
				jwtToken(t, testSecretKeyBytes, time.Now().Add(1*time.Hour)),
			),
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims {
				return jwt.MapClaims{
					"sub": "test",
					"iat": claims["iat"],
					"exp": claims["exp"],
				}
			},
			expectedErr: nil,
		},
		{
			description:    "ShouldReturnErrNotContainTokenInHeader",
			argAccessToken: "",
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims { return nil },
			expectedErr:    auth.ErrNotContainTokenInHeader,
		},
		{
			description:    "ShouldReturnErrInvalidToken",
			argAccessToken: "dfadfasdfasdfasdfasdfsdfsdf",
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims { return nil },
			expectedErr:    auth.ErrInvalidToken,
		},
		{
			description:    "ShouldReturnErrInvalidTokenType",
			argAccessToken: "JWT fasdfasdfasdfasdfasdfsdfasdf",
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims { return nil },
			expectedErr:    auth.ErrInvalidTokenType,
		},
		{
			description:    "ShouldReturnErrInvalidToken",
			argAccessToken: "Bearer fasdfasdfasdfasdfasdfjklasdf",
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims { return nil },
			expectedErr:    auth.ErrInvalidToken,
		},
		{
			description:       "ShouldErrTokenExpired",
			argSecretKeyBytes: testSecretKeyBytes,
			argAccessToken: fmt.Sprintf("Bearer %s",
				jwtToken(t, testSecretKeyBytes, time.Unix(10, 0)),
			),
			expectedClaims: func(claims jwt.MapClaims) jwt.MapClaims { return nil },
			expectedErr:    auth.ErrTokenExpired,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actualClaims, actualErr := verifyAccessToken(tc.argSecretKeyBytes, tc.argAccessToken)

			assert.Equal(t, tc.expectedClaims(actualClaims), actualClaims)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func jwtToken(t *testing.T, secretKeyBytes []byte, expiresAt time.Time) string {
	claims := &jwt.StandardClaims{
		Subject:   "test",
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  time.Unix(1000, 0000).Unix(),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString(secretKeyBytes)

	assert.NoError(t, err)

	return tokenString
}
