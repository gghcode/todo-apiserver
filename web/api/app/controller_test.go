package app_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ControllerUnit struct {
	suite.Suite

	fakeAppService *fake.AppService

	router *gin.Engine
}

func TestCommonControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.New()
	suite.fakeAppService = fake.NewAppService()

	c := app.NewController(suite.fakeAppService)
	c.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestVersion() {
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
				"api/version",
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			actual := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)

			suite.Equal(tc.expectedResponse, actual)
		})
	}
}

func (suite *ControllerUnit) TestHealthy() {
	testCases := []struct {
		description    string
		expectedStatus int
	}{
		{
			description:    "ShouldBeOK",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"api/healthy",
				nil,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}
