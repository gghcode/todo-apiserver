package repository_test

import (
	"testing"

	// "github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type userRepositoryIntegration struct {
	suite.Suite

	dbConn    db.GormConnection
	repo      user.Repository
	dbCleanup func()

	testUsers []user.User
}

func TestUserRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(userRepositoryIntegration))
}

func (suite *userRepositoryIntegration) SetupTest() {
	cfg, err := config.NewViperBuilder().
		BindEnvs("TEST").
		Build()

	suite.NoError(err)

	suite.dbConn, err = db.NewPostgresConn(cfg)
	suite.dbConn.DB().LogMode(false)
	suite.dbCleanup = testutil.DbCleanupFunc(suite.dbConn.DB())
	suite.repo = repository.NewUserRepository(suite.dbConn)

	suite.testUsers = []user.User{
		{UserName: "fakeUser1", PasswordHash: []byte("password")},
		{UserName: "fakeUser2", PasswordHash: []byte("password")},
	}

	for i := range suite.testUsers {
		suite.dbConn.DB().Create(&suite.testUsers[i])
	}
}

func (suite *userRepositoryIntegration) TearDownTest() {
	suite.dbCleanup()
	suite.dbConn.Close()
}

func (suite *userRepositoryIntegration) TestCreateUser() {
	testCases := []struct {
		description string
		argUser     user.User
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldCreateUser",
			argUser:     user.User{UserName: "testuser", PasswordHash: []byte("password")},
			expected:    user.User{UserName: "testuser", PasswordHash: []byte("password")},
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrAlreadyExistUser",
			argUser:     user.User{UserName: suite.testUsers[0].UserName},
			expected:    user.User{},
			expectedErr: user.ErrAlreadyExistUser,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.CreateUser(tc.argUser)

			suite.Equal(tc.expected.UserName, actual.UserName)
			suite.Equal(tc.expected.PasswordHash, actual.PasswordHash)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryIntegration) TestAllUsers() {
	testCases := []struct {
		description string
		expected    []user.User
		expectedErr error
	}{
		{
			description: "ShouldGetAllUsers",
			expected:    suite.testUsers,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.AllUsers()

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryIntegration) TestUserByID() {
	testCases := []struct {
		description string
		argUserID   int64
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserID:   suite.testUsers[0].ID,
			expected:    suite.testUsers[0],
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrNotFoundUser",
			argUserID:   -1,
			expected:    user.User{},
			expectedErr: user.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.UserByID(tc.argUserID)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryIntegration) TestUserByUserName() {
	testCases := []struct {
		description string
		argUserName string
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserName: suite.testUsers[0].UserName,
			expected:    suite.testUsers[0],
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrNotFoundUser",
			argUserName: "NOT_EXIST_USER",
			expected:    user.User{},
			expectedErr: user.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.UserByUserName(tc.argUserName)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryIntegration) TestUpdateUserByID() {
	testCases := []struct {
		description string
		argUser     user.User
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldUpdateUser",
			argUser: user.User{
				ID:           suite.testUsers[0].ID,
				PasswordHash: suite.testUsers[0].PasswordHash,
				CreatedAt:    suite.testUsers[0].CreatedAt,
				UserName:     "updateUserName",
			},
			expected: user.User{
				ID:           suite.testUsers[0].ID,
				PasswordHash: suite.testUsers[0].PasswordHash,
				CreatedAt:    suite.testUsers[0].CreatedAt,
				UserName:     "updateUserName",
			},
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrNotFoundUser",
			argUser:     user.User{ID: -1},
			expected:    user.User{},
			expectedErr: user.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.UpdateUserByID(tc.argUser)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *userRepositoryIntegration) TestRemoveUserByID() {
	testCases := []struct {
		description string
		argUserID   int64
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldRemoveUser",
			argUserID:   suite.testUsers[0].ID,
			expected:    suite.testUsers[0],
			expectedErr: nil,
		},
		{
			description: "ShouldErrNotFoundUser",
			argUserID:   -1,
			expected:    user.User{},
			expectedErr: user.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.RemoveUserByID(tc.argUserID)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
