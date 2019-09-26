package user_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/stretchr/testify/suite"
)

type RepositoryIntegration struct {
	suite.Suite

	postgresConn *db.PostgresConn
	repo         user.Repository
	dbCleanup    func()

	testUsers []user.User
}

func TestUserRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(RepositoryIntegration))
}

func (suite *RepositoryIntegration) SetupTest() {
	cfg, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	suite.NoError(err)

	suite.postgresConn, err = db.NewPostgresConn(cfg)
	suite.postgresConn.DB().LogMode(false)
	suite.dbCleanup = testutil.DbCleanupFunc(suite.postgresConn.DB())
	suite.repo = user.NewRepository(suite.postgresConn)

	suite.testUsers = []user.User{
		{UserName: "fakeUser1", PasswordHash: []byte("password")},
		{UserName: "fakeUser2", PasswordHash: []byte("password")},
	}

	for i := range suite.testUsers {
		suite.postgresConn.DB().Create(&suite.testUsers[i])
	}
}

func (suite *RepositoryIntegration) TearDownTest() {
	suite.dbCleanup()
	suite.postgresConn.Close()
}

func (suite *RepositoryIntegration) TestCreateUser() {
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
			expected:    user.EmptyUser,
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

func (suite *RepositoryIntegration) TestAllUsers() {
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

func (suite *RepositoryIntegration) TestUserByID() {
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
			expected:    user.EmptyUser,
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

func (suite *RepositoryIntegration) TestUserByUserName() {
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
			expected:    user.EmptyUser,
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

func (suite *RepositoryIntegration) TestUpdateUserByID() {
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
			expected:    user.EmptyUser,
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

func (suite *RepositoryIntegration) TestRemoveUserByID() {
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
			expected:    user.EmptyUser,
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
