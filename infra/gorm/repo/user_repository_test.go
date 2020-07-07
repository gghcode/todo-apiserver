package repo_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/entity"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/model"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/repo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type userRepositoryIntegrationTestSuite struct {
	suite.Suite

	repo      user.Repository
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

	dbConn, _, err := gorm.NewPostgresConnection(cfg)
	suite.NoError(err)

	testutil.SetupDBSandbox(suite.T(), dbConn.DB())

	suite.repo = repo.NewUserRepository(dbConn)
	suite.testUsers = []model.User{
		{UserName: "fakeUser1", PasswordHash: []byte("password")},
		{UserName: "fakeUser2", PasswordHash: []byte("password")},
	}

	for i := range suite.testUsers {
		suite.NoError(dbConn.DB().Create(&suite.testUsers[i]).Error)
	}
}

func (suite *userRepositoryIntegrationTestSuite) TestCreateUser() {
	testCases := []struct {
		description string
		argUser     entity.User
		expected    entity.User
		expectedErr error
	}{
		{
			description: "ShouldCreateUser",
			argUser:     entity.User{UserName: "testuser", PasswordHash: []byte("password")},
			expected:    entity.User{UserName: "testuser", PasswordHash: []byte("password")},
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrAlreadyExistUser",
			argUser:     entity.User{UserName: suite.testUsers[0].UserName},
			expected:    entity.User{},
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
		expected    entity.User
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
			expected:    entity.User{},
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
		expected    entity.User
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
			expected:    entity.User{},
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
		argUser     entity.User
		expected    entity.User
		expectedErr error
	}{
		{
			description: "ShouldUpdateUser",
			argUser: entity.User{
				ID:           suite.testUsers[0].ID,
				PasswordHash: suite.testUsers[0].PasswordHash,
				CreatedAt:    suite.testUsers[0].CreatedAt,
				UserName:     "updateUserName",
			},
			expected: entity.User{
				ID:           suite.testUsers[0].ID,
				PasswordHash: suite.testUsers[0].PasswordHash,
				CreatedAt:    suite.testUsers[0].CreatedAt,
				UserName:     "updateUserName",
			},
			expectedErr: nil,
		},
		{
			description: "ShouldBeErrNotFoundUser",
			argUser:     entity.User{ID: -1},
			expected:    entity.User{},
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
		expected    entity.User
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
			expected:    entity.User{},
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
