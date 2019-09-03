package auth_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/auth"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/internal/testutil/fake"
)

type ServiceUnit struct {
	suite.Suite

	cfg      config.Configuration
	jwtParam auth.JWTParam

	fakeUserRepo fake.UserRepository
	fakePassport fake.Passport
	service      auth.Service
}

func TestAuthServiceUnit(t *testing.T) {
	suite.Run(t, new(ServiceUnit))
}

func (suite *ServiceUnit) SetupTest() {
	suite.cfg = config.Configuration{
		Jwt: config.JwtConfig{
			SecretKey:           "testkey",
			AccessExpiresInSec:  3600,
			RefreshExpiresInSec: 7200,
		},
	}

	suite.jwtParam = auth.JWTParam{
		SecretKeyBytes:      []byte(suite.cfg.Jwt.SecretKey),
		AccessExpiresInSec:  time.Duration(suite.cfg.Jwt.AccessExpiresInSec),
		RefreshExpiresInSec: time.Duration(suite.cfg.Jwt.RefreshExpiresInSec),
	}

	suite.fakeUserRepo = fake.UserRepository{}
	suite.fakePassport = fake.Passport{}
	suite.service = auth.NewService(
		suite.cfg,
		&suite.fakePassport,
		&suite.fakeUserRepo,
		fakeCreateAccessToken,
		fakeCreateRefreshToken,
	)
}

func stubCreateAccessToken(p auth.JWTParam, userID int64) string {
	token, _ := fakeCreateAccessToken(p, userID)
	return token
}

func fakeCreateAccessToken(jwtParam auth.JWTParam, userID int64) (string, error) {
	return "access_token", nil
}

func stubCreateRefreshToken(p auth.JWTParam, userID int64) string {
	token, _ := fakeCreateRefreshToken(p, userID)
	return token
}

func fakeCreateRefreshToken(jwtParam auth.JWTParam, userID int64) (string, error) {
	return "refresh_token", nil
}

func (suite *ServiceUnit) TestCreateAccessToken() {
	testCases := []struct {
		description string
		argJwtParam auth.JWTParam
		argUserID   int64
	}{
		{
			description: "ShouldBeEqualSubject",
			argUserID:   100,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			token, err := auth.CreateAccessToken(tc.argJwtParam, tc.argUserID)

			suite.NoError(err)

			claims, err := auth.ExtractTokenClaims(tc.argJwtParam, token)

			suite.NoError(err)

			expected := strconv.FormatInt(tc.argUserID, 10)
			actual := claims["sub"]

			suite.Equal(expected, actual)
		})
	}
}

func (suite *ServiceUnit) TestIssueToken() {
	fakeUser := user.User{
		ID:           100,
		UserName:     "testuser",
		PasswordHash: []byte("testtest"),
	}

	testCases := []struct {
		description string
		argReq      auth.LoginRequest
		stubUser    user.User
		stubErr     error
		stubValid   bool
		expected    auth.TokenResponse
		expectedErr error
	}{
		{
			description: "ShouldIssueToken",
			stubUser:    fakeUser,
			stubErr:     nil,
			stubValid:   true,
			expected: auth.TokenResponse{
				Type:         "Bearer",
				AccessToken:  stubCreateAccessToken(suite.jwtParam, fakeUser.ID),
				RefreshToken: stubCreateRefreshToken(suite.jwtParam, fakeUser.ID),
				ExpiresIn:    suite.cfg.Jwt.AccessExpiresInSec,
			},
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrInvalidCredentialWhenNotFoundUser",
			argReq: auth.LoginRequest{
				Username: "NOT_EXISTS_USER",
				Password: "testtest",
			},
			stubUser:    user.EmptyUser,
			stubErr:     user.ErrUserNotFound,
			expected:    auth.TokenResponse{},
			expectedErr: auth.ErrInvalidCredential,
		},
		{
			description: "ShouldBeErrInvalidCredentialWhenWrongPassword",
			argReq: auth.LoginRequest{
				Username: "test",
				Password: "WRONG_PASSWORD",
			},
			stubUser:    fakeUser,
			stubErr:     nil,
			stubValid:   false,
			expected:    auth.TokenResponse{},
			expectedErr: auth.ErrInvalidCredential,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakePassport.
				On("IsValidPassword", tc.argReq.Password, tc.stubUser.PasswordHash).
				Once().
				Return(tc.stubValid)

			suite.fakeUserRepo.
				On("UserByUserName", tc.argReq.Username).
				Once().
				Return(tc.stubUser, tc.stubErr)

			actual, actualErr := suite.service.IssueToken(tc.argReq)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
