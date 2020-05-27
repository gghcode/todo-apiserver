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
	gin.SetMode(gin.TestMode)

	suite.fakeTodoService = fake.NewTodoService()
	suite.fakeUserIDFactory = &fake.MockUserID{}
	suite.router = gin.New()
	suite.router.Use(gin.HandlerFunc(middleware.NewAccessTokenHandler(
		fake.NewAccessTokenHandlerFactory(suite.fakeUserIDFactory),
	)))

	c := webTodo.NewController(suite.fakeTodoService)
	c.RegisterRoutes(suite.router)
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
				Title:      "title",
				Contents:   "contents",
				AssignorID: 10,
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
				ID:         "id",
				Title:      "test",
				Contents:   "test",
				AssignorID: 10,
			},
			stubErr:        nil,
			expectedStatus: http.StatusCreated,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				map[string]interface{}{
					"id":          "id",
					"title":       "test",
					"contents":    "test",
					"assignor_id": 10,
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
				api.MakeErrorResponseDTO(api.NewUnmarshalError("title", "string")),
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
		{
			description: "ShouldReturnServerInternalError",
			req: todo.AddTodoRequest{
				Title:      "title",
				Contents:   "contents",
				AssignorID: 10,
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
			stubTodoRes:    todo.TodoResponse{},
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON:   testutil.JSONStringFromInterface(suite.T(), api.MakeErrorResponseDTO(fake.ErrStub)),
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
				"/api/todos",
				tc.reqPayload(tc.req),
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)

			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnitTestSuite) TestTodos() {
	testCases := []struct {
		description    string
		stubAuthUserID int64
		stubTodoResArr []todo.TodoResponse
		stubTodoResErr error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldReturnTodos",
			stubAuthUserID: 10,
			stubTodoResArr: []todo.TodoResponse{
				todo.TodoResponse{
					ID:         "TODO_FAKE_ID_1",
					Title:      "title",
					Contents:   "contents",
					AssignorID: 10,
				},
				todo.TodoResponse{
					ID:         "TODO_FAKE_ID_2",
					Title:      "title 2",
					Contents:   "contents 2",
					AssignorID: 10,
				},
				todo.TodoResponse{
					ID:         "TODO_FAKE_ID_3",
					Title:      "title 3",
					Contents:   "contents 3",
					AssignorID: 10,
				},
			},
			stubTodoResErr: nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(), []map[string]interface{}{
				map[string]interface{}{
					"id":          "TODO_FAKE_ID_1",
					"title":       "title",
					"contents":    "contents",
					"assignor_id": 10,
				},
				map[string]interface{}{
					"id":          "TODO_FAKE_ID_2",
					"title":       "title 2",
					"contents":    "contents 2",
					"assignor_id": 10,
				},
				map[string]interface{}{
					"id":          "TODO_FAKE_ID_3",
					"title":       "title 3",
					"contents":    "contents 3",
					"assignor_id": 10,
				},
			}),
		},
		{
			description:    "ShouldReturnServerInternalError",
			stubAuthUserID: 10,
			stubTodoResArr: []todo.TodoResponse{},
			stubTodoResErr: fake.ErrStub,
			expectedStatus: http.StatusInternalServerError,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				api.MakeErrorResponseDTO(fake.ErrStub),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeUserIDFactory.
				On("UserID").
				Once().
				Return(tc.stubAuthUserID)

			suite.fakeTodoService.
				On("GetTodosByUserID", tc.stubAuthUserID).
				Once().
				Return(tc.stubTodoResArr, tc.stubTodoResErr)

			actual := testutil.Response(suite.T(), suite.router, "GET", "/api/todos", nil)
			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)
			suite.JSONEq(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnitTestSuite) TestRemoveTodoByTodoID() {
	testCases := []struct {
		description       string
		queryTodoID       string
		stubAuthUserID    int64
		stubRemoveTodoErr error
		expectedStatus    int
		expectedJSON      string
	}{
		{
			description:       "ShouldRemoveTodo",
			queryTodoID:       "abcd",
			stubAuthUserID:    10,
			stubRemoveTodoErr: nil,
			expectedStatus:    http.StatusNoContent,
			expectedJSON:      "",
		},
		{
			description:       "ShouldReturnNotFoundTodo",
			queryTodoID:       "abcd",
			stubAuthUserID:    10,
			stubRemoveTodoErr: todo.ErrNotFoundTodo,
			expectedStatus:    http.StatusNotFound,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				api.MakeErrorResponseDTO(todo.ErrNotFoundTodo),
			),
		},
		{
			description:       "ShouldReturnServerInternalError",
			queryTodoID:       "abcd",
			stubAuthUserID:    10,
			stubRemoveTodoErr: fake.ErrStub,
			expectedStatus:    http.StatusInternalServerError,
			expectedJSON: testutil.JSONStringFromInterface(suite.T(),
				api.MakeErrorResponseDTO(fake.ErrStub),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.fakeUserIDFactory.
				On("UserID").
				Once().
				Return(tc.stubAuthUserID)

			suite.fakeTodoService.
				On("RemoveTodo", tc.queryTodoID).
				Once().
				Return(tc.stubRemoveTodoErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"DELETE",
				"/api/todos/"+tc.queryTodoID,
				nil,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(suite.T(), actual.Body)
			if tc.expectedJSON == "" {
				suite.Equal(tc.expectedJSON, actualJSON)
			} else {
				suite.JSONEq(tc.expectedJSON, actualJSON)
			}
		})
	}
}
