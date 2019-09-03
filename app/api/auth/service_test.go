package auth_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/auth"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/internal/testutil/fake"
)

type ServiceUnit struct {
	suite.Suite

	cfg          config.Configuration
	fakeUserRepo fake.UserRepository
	fakePassport fake.Passport
	service      auth.Service
}

func TestAuthServiceUnit(t *testing.T) {
	suite.Run(t, new(ServiceUnit))
}

func (suite *ServiceUnit) SetupTest() {
	suite.cfg = config.Configuration{}

	suite.fakeUserRepo = fake.UserRepository{}
	suite.fakePassport = fake.Passport{}
	suite.service = auth.NewService(
		suite.cfg,
		&suite.fakePassport,
		&suite.fakeUserRepo,
	)
}

func (suite *ServiceUnit) TestIssueToken() {
	fakeUser := user.User{
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
			expected:    auth.TokenResponse{},
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
