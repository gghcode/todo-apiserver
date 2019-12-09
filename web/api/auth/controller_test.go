package auth_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerUnitTestSuite struct {
	suite.Suite

	router          *gin.Engine
	fakeAuthService *fake.AuthService
}

func TestAuthControllerUnitTests(t *testing.T) {
	suite.Run(t, new(ControllerUnitTestSuite))
}

func (suite *ControllerUnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.router = gin.New()
	suite.fakeAuthService = fake.NewAuthService()

	c := webAuth.NewController(suite.fakeAuthService)
	c.RegisterRoutes(suite.router)
}

func (suite *ControllerUnitTestSuite) TestIssueToken() {
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
			expectedJSON:   `{"error":{"message":"username: cannot be blank."}}`,
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
				"api/auth/token",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnitTestSuite) TestRefreshToken() {
	fakeTokenResponse := auth.TokenResponse{
		Type:        "Bearer",
		AccessToken: "abadfasdf",
		ExpiresIn:   3600,
	}

	testCases := []struct {
		description    string
		reqPayload     io.Reader
		stubToken      auth.TokenResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldRefreshToken",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				auth.AccessTokenByRefreshRequest{Token: "fasdfasdf"},
			),
			stubToken:      fakeTokenResponse,
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				fakeTokenResponse,
			),
		},
		{
			description: "ShouldReturnBadRequestWhenEmptyToken",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				auth.AccessTokenByRefreshRequest{Token: ""},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"token: cannot be blank."}}`,
		},
		{
			description: "ShouldReturnBadRequestWhenNotContainToken",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				map[string]interface{}{},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"token: cannot be blank."}}`,
		},
		{
			description: "ShouldReturnBadRequestWhenTokenTypeInteger",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				map[string]interface{}{
					"token": 3,
				},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponse(
					api.NewUnmarshalError("token", "string"),
				),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeAuthService.
				On("RefreshToken", mock.Anything).
				Once().
				Return(tc.stubToken, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/auth/refresh",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}
