package common_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api/common"
	"github.com/gghcode/apas-todo-apiserver/app/loader"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/suite"
)

type ControllerUnit struct {
	suite.Suite

	router     *gin.Engine
	controller *common.Controller
}

func TestCommonControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.New()

	versionLoader := loader.NewVersionLoader(afero.NewOsFs())

	suite.controller = common.NewController(versionLoader)
	suite.controller.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestVersion() {
	testCases := []struct {
		description      string
		expectedStatus   int
		expectedResponse string
	}{
		{
			description:      "ShouldBeDevVersion",
			expectedStatus:   http.StatusOK,
			expectedResponse: "dev version",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
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
