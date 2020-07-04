package repository_test

import (
	"testing"
	"time"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/stretchr/testify/suite"
)

type RepositoryIntegration struct {
	suite.Suite

	redisConn db.RedisConnection
	tokenRepo auth.TokenRepository

	cleanup func()
}

func TestTokenRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(RepositoryIntegration))
}

func (suite *RepositoryIntegration) SetupTest() {
	cfg, err := config.FromEnvs()
	suite.NoError(err)

	suite.redisConn, suite.cleanup = db.NewRedisConn(cfg)
	suite.tokenRepo = repository.NewRedisTokenRepository(suite.redisConn)
}

func (suite *RepositoryIntegration) TearDownTest() {
	suite.cleanup()
}

func (suite *RepositoryIntegration) TestSaveRefreshToken() {
	testCases := []struct {
		description     string
		argUserID       int64
		argRefreshToken string
		argExpireIn     time.Duration
		expected        error
	}{
		{
			description:     "ShouldSaveToken",
			argUserID:       100,
			argRefreshToken: "refreshtoken",
			expected:        nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual := suite.tokenRepo.SaveRefreshToken(
				tc.argUserID,
				tc.argRefreshToken,
				tc.argExpireIn,
			)

			suite.Equal(tc.expected, actual)
		})
	}
}

func (suite *RepositoryIntegration) TestUserIDByRefreshToken() {
	fakeUserID := int64(100)
	fakeTokenString := "debug token"

	suite.tokenRepo.SaveRefreshToken(
		fakeUserID,
		fakeTokenString,
		time.Duration(time.Now().AddDate(0, 0, 1).Unix()),
	)

	testCases := []struct {
		description     string
		argRefreshToken string
		expected        int64
		expectedErr     error
	}{
		{
			description:     "ShouldFetchRefreshToken",
			argRefreshToken: fakeTokenString,
			expected:        fakeUserID,
			expectedErr:     nil,
		},
		{
			description:     "ShouldReturnErrNotStoredToken",
			argRefreshToken: "",
			expected:        -1,
			expectedErr:     auth.ErrNotStoredToken,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.tokenRepo.UserIDByRefreshToken(tc.argRefreshToken)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
