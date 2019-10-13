package user_test

import (
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/app/infra"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ControllerUnit struct {
	suite.Suite

	router         *gin.Engine
	controller     *user.Controller
	userIDFactory  *fake.MockUserID
	userRepository *fake.UserRepository
}

func TestUserControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.userIDFactory = &fake.MockUserID{}
	suite.router = gin.New()
	suite.router.Use(fake.AddJwtAuthHandler(suite.userIDFactory))
	suite.userRepository = &fake.UserRepository{}

	suite.controller = user.NewController(suite.userRepository, infra.NewPassport(0))
	suite.controller.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestCreateUser() {
	testCases := []struct {
		description     string
		reqPayload      io.Reader
		stubCreatedUser user.User
		stubErr         error
		expectedStatus  int
		expectedJSON    string
	}{
		{
			description: "ShouldCreateUser",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{
					UserName: "testUser",
					Password: "password",
				},
			),
			stubCreatedUser: user.User{
				ID:           1,
				UserName:     "testUser",
				PasswordHash: []byte("password"),
				CreatedAt:    1000,
			},
			stubErr:        nil,
			expectedStatus: http.StatusCreated,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), user.UserResponse{
				ID:        1,
				UserName:  "testUser",
				CreatedAt: time.Unix(1000, 0),
			}),
		},
		{
			description:    "ShouldBeBadRequestWhenNotContainPayload",
			reqPayload:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   "-",
		},
		{
			description: "ShouldBeBadRequestWhenEmptyUserName",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{Password: "password"},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   "-",
		},
		{
			description: "ShouldBeBadRequestWhenEmptyPassword",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				user.CreateUserRequest{UserName: "testuser"},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   "-",
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
			expectedJSON:   "-",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.userRepository.
				On("CreateUser", mock.Anything).
				Return(tc.stubCreatedUser, tc.stubErr)

			actualRes := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/users",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actualRes.StatusCode)

			if tc.expectedJSON != "-" {
				actualJSON := testutil.JSONStringFromResBody(suite.T(), actualRes.Body)
				suite.Equal(tc.expectedJSON, actualJSON)
			}
		})
	}
}

func (suite *ControllerUnit) TestAuthenticatedUser() {
	testUser := user.User{
		ID:        100,
		UserName:  "testuser",
		CreatedAt: time.Now().Unix(),
	}

	testCases := []struct {
		description    string
		stubUserID     int64
		stubUser       user.User
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldReturnUser",
			stubUserID:     testUser.ID,
			stubUser:       testUser,
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), testUser.Response()),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.userIDFactory.
				On("UserID").
				Return(tc.stubUserID)

			suite.userRepository.
				On("UserByID", tc.stubUserID).
				Return(tc.stubUser, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"api/user",
				nil,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnit) TestUserByName() {
	testUser := user.User{
		ID:       10,
		UserName: "testUser",
	}

	testCases := []struct {
		description    string
		argUserName    string
		stubUser       user.User
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldFetchuser",
			argUserName:    testUser.UserName,
			stubUser:       testUser,
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				testUser.Response(),
			),
		},
		{
			description:    "ShouldBeUserNotExists",
			argUserName:    "NOT_EXISTS_USER",
			stubUser:       user.EmptyUser,
			stubErr:        user.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.NewErrRes(user.ErrUserNotFound),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.userRepository.
				On("UserByUserName", tc.argUserName).
				Return(tc.stubUser, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"api/users/"+tc.argUserName,
				nil,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}
