package user_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/internal/testutil"
)

type ControllerUnit struct {
	suite.Suite

	router     *gin.Engine
	controller *user.Controller
}

func TestUserControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.New()

	suite.controller = user.NewController()
	suite.controller.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestCreateUser() {
	testCases := []struct {
		description    string
		reqPayload     io.Reader
		expectedStatus int
	}{
		{
			description: "ShouldCreateUser",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{
					UserName: "testuser",
					Password: "password",
				},
			),
			expectedStatus: http.StatusCreated,
		},
		{
			description:    "ShouldBeBadRequestWhenNotContainPayload",
			reqPayload:     nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "ShouldBeBadRequestWhenEmptyUserName",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{Password: "password"},
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "ShouldBeBadRequestWhenEmptyPassword",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{UserName: "testuser"},
			),
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "ShouldBeBadRequestWhenInvalidPassword",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{
					UserName: "testuser",
					Password: "1234",
				},
			),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualRes := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"/users",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)
		})
	}
}
