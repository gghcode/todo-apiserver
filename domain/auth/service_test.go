package auth_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceUnit struct {
	suite.Suite

	cfg config.Configuration

	fakeTokenRepo fake.TokenRepository
	fakeUserRepo  fake.UserRepository
	fakePassport  fake.Passport
	service       auth.UsecaseInteractor
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

	suite.fakeTokenRepo = fake.TokenRepository{}
	suite.fakeUserRepo = fake.UserRepository{}
	suite.fakePassport = fake.Passport{}
	suite.service = auth.NewService(
		suite.cfg,
		&suite.fakePassport,
		&suite.fakeTokenRepo,
		&suite.fakeUserRepo,
		fakeCreateAccessTokenFactory,
		fakeCreateRefreshTokenFactory,
	)
}

func stubCreateAccessToken(userID int64) string {
	handler := fakeCreateAccessTokenFactory(config.DefaultConfig())
	token, _ := handler(userID)
	return token
}

func fakeCreateAccessTokenFactory(cfg config.Configuration) auth.CreateAccessTokenHandler {
	return func(userID int64) (string, error) {
		return "access_token", nil
	}
}

func stubCreateRefreshToken(tokenRepo auth.TokenRepository, userID int64) string {
	handler := fakeCreateRefreshTokenFactory(config.DefaultConfig(), tokenRepo)
	token, _ := handler(userID)
	return token
}

func fakeCreateRefreshTokenFactory(cfg config.Configuration, tokenRepo auth.TokenRepository) auth.CreateRefreshTokenHandler {
	return func(userID int64) (string, error) {
		return "refresh_token", nil
	}
}

// func (suite *ServiceUnit) TestCreateRefreshToken() {
// 	argJwtParam := suite.jwtParam
// 	argUserID := int64(100)
// 	expectedSub := strconv.FormatInt(argUserID, 10)

// 	suite.fakeTokenRepo.On("SaveRefreshToken",
// 		argUserID,
// 		mock.Anything,
// 		argJwtParam.RefreshExpiresInSec,
// 	).Return(nil)

// 	createRefreshTokenHandler := auth.CreateRefreshTokenFactory(argJwtParam, &suite.fakeTokenRepo)
// 	token, err := createRefreshTokenHandler(argUserID)
// 	suite.NoError(err)

// 	claims, err := auth.ExtractTokenClaims(argJwtParam, token)
// 	suite.NoError(err)

// 	suite.Equal(expectedSub, claims["sub"])
// }

// func (suite *ServiceUnit) TestCreateAccessToken() {
// 	testCases := []struct {
// 		description string
// 		argJwtParam auth.JwtParam
// 		argUserID   int64
// 	}{
// 		{
// 			description: "ShouldBeEqualSubject",
// 			argUserID:   100,
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.Run(tc.description, func() {
// 			createAccessTokenHandler := auth.CreateAccessTokenFactory(tc.argJwtParam)
// 			token, err := createAccessTokenHandler(tc.argUserID)

// 			suite.NoError(err)

// 			claims, err := auth.ExtractTokenClaims(tc.argJwtParam, token)

// 			suite.NoError(err)

// 			expected := strconv.FormatInt(tc.argUserID, 10)
// 			actual := claims["sub"]

// 			suite.Equal(expected, actual)
// 		})
// 	}
// }

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
				AccessToken:  stubCreateAccessToken(fakeUser.ID),
				RefreshToken: stubCreateRefreshToken(&suite.fakeTokenRepo, fakeUser.ID),
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
			stubUser:    user.User{},
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

			suite.fakeTokenRepo.On("SaveRefreshToken",
				tc.stubUser.ID,
				mock.Anything,
				mock.Anything,
			).Return(nil)

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

func (suite *ServiceUnit) TestRefreshToken() {
	testCases := []struct {
		description string
		argReq      auth.AccessTokenByRefreshRequest
		stubUserID  int64
		stubErr     error
		expected    auth.TokenResponse
		expectedErr error
	}{
		{
			description: "ShouldRefreshToken",
			stubUserID:  100,
			stubErr:     nil,
			expected: auth.TokenResponse{
				Type:         "Bearer",
				AccessToken:  stubCreateAccessToken(10),
				RefreshToken: "",
				ExpiresIn:    3600,
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnNotStoredErr",
			stubUserID:  0,
			stubErr:     auth.ErrNotStoredToken,
			expected:    auth.TokenResponse{},
			expectedErr: auth.ErrNotStoredToken,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeTokenRepo.
				On("UserIDByRefreshToken", mock.Anything).
				Once().
				Return(tc.stubUserID, tc.stubErr)

			actual, actualErr := suite.service.RefreshToken(tc.argReq)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
