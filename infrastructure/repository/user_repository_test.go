package repository_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/model"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type userRepositoryIntegrationTestSuite struct {
	suite.Suite

	repo      user.Repository
	dbCleanup func()

	testUsers []model.User
}

func TestUserRepositoryIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(userRepositoryIntegrationTestSuite))
}

func (suite *userRepositoryIntegrationTestSuite) SetupTest() {
	cfg, err := config.FromEnvs()
	suite.NoError(err)

	dbConn, _, err := db.NewPostgresConn(cfg)
	suite.NoError(err)

	suite.dbCleanup = testutil.DbCleanupFunc(dbConn.DB())
	suite.repo = repository.NewUserRepository(dbConn)
	suite.testUsers = []model.User{
		{UserName: "fakeUser1", PasswordHash: []byte("password")},
		{UserName: "fakeUser2", PasswordHash: []byte("password")},
	}

	for i := range suite.testUsers {
		suite.NoError(dbConn.DB().Create(&suite.testUsers[i]).Error)
	}
}

func (suite *userRepositoryIntegrationTestSuite) TearDownTest() {
	suite.dbCleanup()
}

func (suite *userRepositoryIntegrationTestSuite) TestCreateUser() {
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

func (suite *userRepositoryIntegrationTestSuite) TestUserByID() {
	testCases := []struct {
		description string
		argUserID   int64
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserID:   suite.testUsers[0].ID,
			expected:    model.ToUserEntity(suite.testUsers[0]),
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

func (suite *userRepositoryIntegrationTestSuite) TestUserByUserName() {
	testCases := []struct {
		description string
		argUserName string
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserName: suite.testUsers[0].UserName,
			expected:    model.ToUserEntity(suite.testUsers[0]),
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

func (suite *userRepositoryIntegrationTestSuite) TestUpdateUserByID() {
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

func (suite *userRepositoryIntegrationTestSuite) TestRemoveUserByID() {
	testCases := []struct {
		description string
		argUserID   int64
		expected    user.User
		expectedErr error
	}{
		{
			description: "ShouldRemoveUser",
			argUserID:   suite.testUsers[0].ID,
			expected:    model.ToUserEntity(suite.testUsers[0]),
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
