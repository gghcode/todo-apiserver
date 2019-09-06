package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/auth"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
)

type RepositoryIntegration struct {
	suite.Suite

	redisConn db.RedisConn
	tokenRepo auth.Repository
}

func TestTokenRepositoryIntegration(t *testing.T) {
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

	suite.redisConn = db.NewRedisConn(cfg)
	suite.tokenRepo = auth.NewRepository(suite.redisConn)
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

func (suite *RepositoryIntegration) TestRefreshToken() {
	fakeUserID := int64(100)
	fakeTokenString := "debug token"

	suite.tokenRepo.SaveRefreshToken(
		fakeUserID,
		fakeTokenString,
		time.Duration(time.Now().AddDate(0, 0, 1).Unix()),
	)

	testCases := []struct {
		description string
		argUserID   int64
		expected    string
		expectedErr error
	}{
		{
			description: "ShouldFetchRefreshToken",
			argUserID:   fakeUserID,
			expected:    fakeTokenString,
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotStoredToken",
			argUserID:   -1,
			expected:    "",
			expectedErr: auth.ErrNotStoredToken,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.tokenRepo.RefreshToken(tc.argUserID)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
