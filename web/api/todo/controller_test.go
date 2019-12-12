package todo_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webTodo "github.com/gghcode/apas-todo-apiserver/web/api/todo"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ControllerUnitTestSuite struct {
	suite.Suite

	router            *gin.Engine
	fakeTodoService   *fake.TodoService
	fakeUserIDFactory *fake.MockUserID
}

func TestTodoControllerUnitTests(t *testing.T) {
	suite.Run(t, new(ControllerUnitTestSuite))
}

func (suite *ControllerUnitTestSuite) SetupTest() {
	suite.fakeTodoService = fake.NewTodoService()
	suite.fakeUserIDFactory = &fake.MockUserID{}
	suite.router = gin.New()
	suite.router.Use(middleware.AddAccessTokenHandler(
		fake.NewAccessTokenHandlerFactory(suite.fakeUserIDFactory),
	))

	c := webTodo.NewController(suite.fakeTodoService)
	c.RegisterRoutes(suite.router)

	gin.SetMode(gin.TestMode)
}

func (suite *ControllerUnitTestSuite) TestAddTodo() {
	testCases := []struct {
		description    string
		req            todo.AddTodoRequest
		reqPayload     func(req todo.AddTodoRequest) io.Reader
		stubAuthUserID int64
		stubTodoRes    todo.TodoResponse
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldAddTodo",
			req: todo.AddTodoRequest{
				Title:    "title",
				Contents: "contents",
			},
			reqPayload: func(req todo.AddTodoRequest) io.Reader {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					map[string]interface{}{
						"title":    req.Title,
						"contents": req.Contents,
					},
				)
			},
			stubAuthUserID: 10,
			stubTodoRes: todo.TodoResponse{
				ID:       "id",
				Title:    "test",
				Contents: "test",
			},
			stubErr:        nil,
			expectedStatus: http.StatusCreated,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				map[string]interface{}{
					"id":       "id",
					"title":    "test",
					"contents": "test",
				},
			),
		},
		{
			description: "ShouldBadRequestWhenInvalidType",
			reqPayload: func(req todo.AddTodoRequest) io.Reader {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					map[string]interface{}{
						"title":    3,
						"contents": "contents",
					},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.MakeErrorResponse(api.NewUnmarshalError("title", "string")),
			),
		},
		{
			description: "ShouldBadRequestWhenNotContainContents",
			reqPayload: func(req todo.AddTodoRequest) io.Reader {
				return testutil.ReqBodyFromInterface(
					suite.T(),
					map[string]interface{}{
						"title": "new title",
					},
				)
			},
			expectedStatus: http.StatusBadRequest,
			expectedJSON:   `{"error":{"message":"contents: cannot be blank."}}`,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeUserIDFactory.
				On("UserID").
				Once().
				Return(tc.stubAuthUserID)

			suite.fakeTodoService.
				On("AddTodo", tc.req).
				Once().
				Return(tc.stubTodoRes, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/todos",
				tc.reqPayload(tc.req),
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)

			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}
