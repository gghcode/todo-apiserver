package auth_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/auth"
	"gitlab.com/gyuhwan/apas-todo-apiserver/internal/testutil"
)

type fakeAuthService struct {
	mock.Mock
}

func (service *fakeAuthService) IssueToken(req auth.LoginRequest) (auth.TokenResponse, error) {
	args := service.Called(req)
	return args.Get(0).(auth.TokenResponse), args.Error(1)
}

type ControllerUnit struct {
	suite.Suite

	router          *gin.Engine
	fakeAuthService *fakeAuthService
}

func TestAuthControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.New()
	suite.fakeAuthService = &fakeAuthService{}

	controller := auth.NewController(suite.fakeAuthService)
	controller.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestGetToken() {
	fakeTokenRes := auth.TokenResponse{
		Type:         "Bearer",
		AccessToken:  "fasdf",
		RefreshToken: "fasdf",
		ExpiresIn:    123,
	}

	testCases := []struct {
		description    string
		reqPayload     io.Reader
		stubTokenRes   auth.TokenResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldGenerateToken",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				auth.LoginRequest{
					Username: "test",
					Password: "testtest",
				},
			),
			stubTokenRes:   fakeTokenRes,
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), fakeTokenRes),
		},
		{
			description: "ShouldBeBadRequestWhenEmptyUsername",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				auth.LoginRequest{
					Username: "",
					Password: "testtest",
				},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.NewErrRes(
					api.Validate(
						auth.LoginRequest{
							Username: "",
							Password: "testtest",
						},
					),
				),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeAuthService.
				On("IssueToken", mock.Anything).
				Once().
				Return(tc.stubTokenRes, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				auth.APIPath+"/token",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}
