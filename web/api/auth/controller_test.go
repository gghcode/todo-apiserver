package auth_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"

	"github.com/gin-gonic/gin"
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

	suite.fakeAuthService = fake.NewAuthService()
	suite.router = gin.New()

	c := webAuth.NewController(suite.fakeAuthService)
	c.RegisterRoutes(suite.router)
}

func (suite *ControllerUnitTestSuite) TestIssueToken() {
	testCases := []struct {
		description    string
		req            auth.LoginRequest
		reqPayload     func(req auth.LoginRequest) *bytes.Buffer
		stubTokenRes   auth.TokenResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldGenerateToken",
			req: auth.LoginRequest{
				Username: "test",
				Password: "testtest",
			},
			reqPayload: func(req auth.LoginRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubTokenRes: auth.TokenResponse{
				Type:         "Bearer",
				AccessToken:  "fasdf",
				RefreshToken: "fasdf",
				ExpiresIn:    123,
			},
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), map[string]interface{}{
				"type":          "Bearer",
				"access_token":  "fasdf",
				"refresh_token": "fasdf",
				"expires_in":    123,
			}),
		},
		{
			description: "ShouldBeBadRequestWhenEmptyUsername",
			reqPayload: func(req auth.LoginRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					auth.LoginRequest{
						Username: "",
						Password: "testtest",
					},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"username: cannot be blank."}}`,
		},
		{
			description: "ShouldBeBadRequestWhenPasswordInteger",
			reqPayload: func(req auth.LoginRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
					"username": "test",
					"password": 34,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponseDTO(
				api.NewUnmarshalError("password", "string"),
			)),
		},
		{
			description: "ShouldBeUnauthorizedWhenInvalidCredential",
			req: auth.LoginRequest{
				Username: "test",
				Password: "testtest",
			},
			reqPayload: func(req auth.LoginRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubTokenRes:   auth.TokenResponse{},
			stubErr:        auth.ErrInvalidCredential,
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponseDTO(
					auth.ErrInvalidCredential,
				),
			),
		},
		{
			description: "ShouldReturnServerInternalError",
			req: auth.LoginRequest{
				Username: "test",
				Password: "testtest",
			},
			reqPayload: func(req auth.LoginRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubTokenRes:   auth.TokenResponse{},
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponseDTO(fake.ErrStub)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeAuthService.
				On("IssueToken", tc.req).
				Once().
				Return(tc.stubTokenRes, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/auth/token",
				tc.reqPayload(tc.req),
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)

			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnitTestSuite) TestRefreshToken() {
	testCases := []struct {
		description    string
		req            auth.AccessTokenByRefreshRequest
		reqPayload     func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer
		stubToken      auth.TokenResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldRefreshToken",
			req:         auth.AccessTokenByRefreshRequest{Token: "fasdfasdf"},
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubToken: auth.TokenResponse{
				Type:        "Bearer",
				AccessToken: "abadfasdf",
				ExpiresIn:   3600,
			},
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				map[string]interface{}{
					"type":         "Bearer",
					"access_token": "abadfasdf",
					"expires_in":   3600,
				},
			),
		},
		{
			description: "ShouldReturnBadRequestWhenEmptyToken",
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					auth.AccessTokenByRefreshRequest{Token: ""},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"token: cannot be blank."}}`,
		},
		{
			description: "ShouldReturnBadRequestWhenNotContainToken",
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					map[string]interface{}{},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"token: cannot be blank."}}`,
		},
		{
			description: "ShouldReturnBadRequestWhenTokenTypeInteger",
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					map[string]interface{}{
						"token": 3,
					},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponseDTO(
					api.NewUnmarshalError("token", "string"),
				),
			),
		},
		{
			description: "ShouldReturnUnauthrizedWhenErrNotStoredToken",
			req:         auth.AccessTokenByRefreshRequest{Token: "abcd"},
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					req,
				)
			},
			stubToken:      auth.TokenResponse{},
			stubErr:        auth.ErrNotStoredToken,
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponseDTO(
					auth.ErrNotStoredToken,
				),
			),
		},
		{
			description: "ShouldReturnServerInternalError",
			req:         auth.AccessTokenByRefreshRequest{Token: "abcdabcd"},
			reqPayload: func(req auth.AccessTokenByRefreshRequest) *bytes.Buffer {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					req,
				)
			},
			stubToken:      auth.TokenResponse{},
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponseDTO(fake.ErrStub)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeAuthService.
				On("RefreshToken", tc.req).
				Once().
				Return(tc.stubToken, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/auth/refresh",
				tc.reqPayload(tc.req),
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)

			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}
