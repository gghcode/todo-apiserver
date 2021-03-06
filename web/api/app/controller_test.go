package app_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/redis"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ControllerIntegrationTestSuite struct {
	suite.Suite

	fakeAppService *fake.AppService
	postgresConn   gorm.Connection
	redisConn      redis.Connection
	cleanup        func()

	router *gin.Engine
}

func TestAppControllerIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(ControllerIntegrationTestSuite))
}

func (suite *ControllerIntegrationTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	cfg, err := config.FromEnvs()
	suite.NoError(err)

	suite.router = gin.New()
	suite.fakeAppService = fake.NewAppService()

	postgresConn, postgresCleanup, err := gorm.NewPostgresConnection(cfg)
	redisConn, redisCleanup := redis.NewConnection(cfg)

	suite.NoError(err)

	suite.postgresConn = postgresConn
	suite.redisConn = redisConn
	suite.cleanup = func() {
		redisCleanup()
		postgresCleanup()
	}

	c := app.NewController(suite.fakeAppService, suite.postgresConn, suite.redisConn)
	c.RegisterRoutes(suite.router)
}

func (suite *ControllerIntegrationTestSuite) TearDownTest() {
	suite.cleanup()
}

func (suite *ControllerIntegrationTestSuite) TestVersion() {
	testCases := []struct {
		description      string
		stubVersion      string
		expectedStatus   int
		expectedResponse string
	}{
		{
			description:      "ShouldBeDevVersion",
			stubVersion:      "test version",
			expectedStatus:   http.StatusOK,
			expectedResponse: "test version",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeAppService.
				On("Version").
				Return(tc.stubVersion)

			actualRes := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"/api/version",
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actual := testutil.StringFromIOReader(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedResponse, actual)
		})
	}
}

func (suite *ControllerIntegrationTestSuite) TestHealthy() {
	testCases := []struct {
		description    string
		beforeFunc     func()
		expectedStatus int
	}{
		{
			description:    "ShouldBeOK",
			beforeFunc:     func() {},
			expectedStatus: http.StatusOK,
		},
		{
			description: "ShouldBeServiceUnavailableWhenDownRedis",
			beforeFunc: func() {
				suite.redisConn.Close()
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			description: "ShouldBeServiceUnavailableWhenDownPostgres",
			beforeFunc: func() {
				suite.postgresConn.Close()
			},
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			tc.beforeFunc()

			actualRes := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"/api/healthy",
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}
