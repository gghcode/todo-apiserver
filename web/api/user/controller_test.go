package user_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webUser "github.com/gghcode/apas-todo-apiserver/web/api/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ControllerUnitTestSuite struct {
	suite.Suite

	router          *gin.Engine
	fakeUserService *fake.UserService
}

func TestUserControllerUnitTests(t *testing.T) {
	suite.Run(t, new(ControllerUnitTestSuite))
}

func (suite *ControllerUnitTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.fakeUserService = fake.NewUserService()
	suite.router = gin.New()

	c := webUser.NewController(suite.fakeUserService)
	c.RegisterRoutes(suite.router)
}

func (suite *ControllerUnitTestSuite) TestCreateUser() {
	testCases := []struct {
		description    string
		req            user.CreateUserRequest
		reqPayload     func(req user.CreateUserRequest) io.Reader
		stubUserRes    user.UserResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldCreateUser",
			req: user.CreateUserRequest{
				UserName: "test",
				Password: "testtest",
			},
			reqPayload: func(req user.CreateUserRequest) io.Reader {
				return testutil.ReqBodyFromInterface(suite.T(),
					map[string]interface{}{
						"username": req.UserName,
						"password": req.Password,
					})
			},
			stubUserRes: user.UserResponse{
				ID:       1,
				UserName: "test",
			},
			stubErr:        nil,
			expectedStatus: http.StatusCreated,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				map[string]interface{}{
					"id":         1,
					"username":   "test",
					"created_at": "0001-01-01T00:00:00Z",
				},
			),
		},
		{
			description: "ShouldBadRequestWhenEmptyUserName",
			req: user.CreateUserRequest{
				UserName: "",
				Password: "testtest",
			},
			reqPayload: func(req user.CreateUserRequest) io.Reader {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"username: cannot be blank."}}`,
		},
		{
			description: "ShouldConflictWhenAlreadyExistUser",
			req: user.CreateUserRequest{
				UserName: "test",
				Password: "testtest",
			},
			reqPayload: func(req user.CreateUserRequest) io.Reader {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubUserRes:    user.UserResponse{},
			stubErr:        user.ErrAlreadyExistUser,
			expectedStatus: http.StatusConflict,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponse(user.ErrAlreadyExistUser),
			),
		},
		{
			description: "ShouldReturnBadRequestWhenPasswordInteger",
			reqPayload: func(req user.CreateUserRequest) io.Reader {
				return testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
					"username": "test",
					"password": 3,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponse(
				api.NewUnmarshalError("password", "string"),
			)),
		},
		{
			description: "ShouldReturnServerInternalError",
			req: user.CreateUserRequest{
				UserName: "test",
				Password: "testtest",
			},
			reqPayload: func(req user.CreateUserRequest) io.Reader {
				return testutil.ReqBodyFromInterface(suite.T(), req)
			},
			stubUserRes:    user.UserResponse{},
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponse(fake.ErrStub)),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeUserService.
				On("CreateUser", tc.req).
				Once().
				Return(tc.stubUserRes, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/users",
				tc.reqPayload(tc.req),
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)

			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}
