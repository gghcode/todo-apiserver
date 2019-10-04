package todo_test

import (
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type fakeTodoRepository struct {
	mock.Mock
}

func (repo *fakeTodoRepository) AddTodo(t todo.Todo) (todo.Todo, error) {
	args := repo.Called(t)
	return args.Get(0).(todo.Todo), args.Error(1)
}

func (repo *fakeTodoRepository) AllTodosByUserID(userID int64) ([]todo.Todo, error) {
	args := repo.Called(userID)
	return args.Get(0).([]todo.Todo), args.Error(1)
}

type ControllerUnit struct {
	suite.Suite

	router        *gin.Engine
	userIDFactory *fake.MockUserID
	controller    *todo.Controller
	todoRepo      *fakeTodoRepository
}

func TestTodoControllerUnit(t *testing.T) {
	suite.Run(t, new(ControllerUnit))
}

func (suite *ControllerUnit) SetupTest() {
	gin.SetMode(gin.TestMode)

	suite.userIDFactory = &fake.MockUserID{}
	suite.router = gin.New()
	suite.router.Use(fake.AddJwtAuthHandler(suite.userIDFactory))

	suite.todoRepo = &fakeTodoRepository{}

	suite.controller = todo.NewController(suite.todoRepo)
	suite.controller.RegisterRoutes(suite.router)
}

func (suite *ControllerUnit) TestAddTodo() {
	fakeTodo := todo.Todo{
		ID:         uuid.NewV4(),
		Title:      "Fake Title",
		Contents:   "fake contents",
		AssignorID: 10,
	}

	testCases := []struct {
		description    string
		reqPayload     io.Reader
		stubTodo       todo.Todo
		stubErr        error
		stubUserID     int64
		expectedStatus int
		expectedJSON   string
	}{
		{
			description: "ShouldAddTodo",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				todo.AddTodoRequest{
					Title:    fakeTodo.Title,
					Contents: fakeTodo.Contents,
				},
			),
			stubTodo:       fakeTodo,
			stubErr:        nil,
			stubUserID:     fakeTodo.AssignorID,
			expectedStatus: http.StatusCreated,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				todo.TodoSerializer{Model: fakeTodo}.Response(),
			),
		},
		{
			description: "ShouldBeBadRequestWhenEmptyTitle",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				todo.AddTodoRequest{Title: "", Contents: "contents"},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.NewErrRes(
					api.Validate(todo.AddTodoRequest{Title: "", Contents: "contents"}),
				),
			),
		},
		{
			description: "ShouldBeBadRequestWhenEmptyContents",
			reqPayload: testutil.ReqBodyFromInterface(
				suite.T(),
				todo.AddTodoRequest{Title: "title", Contents: ""},
			),
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.NewErrRes(
					api.Validate(todo.AddTodoRequest{Title: "title", Contents: ""}),
				),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			suite.userIDFactory.
				On("UserID").
				Return(tc.stubUserID)

			suite.todoRepo.
				On("AddTodo", mock.Anything).
				Once().
				Return(tc.stubTodo, tc.stubErr)

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"api/todos",
				tc.reqPayload,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}

func (suite *ControllerUnit) TestAllTodos() {
	fakeTodos := []todo.Todo{
		{
			ID:       uuid.NewV4(),
			Title:    "fake title",
			Contents: "fake contents",
		},
		{
			ID:       uuid.NewV4(),
			Title:    "fake title2",
			Contents: "fake contents2",
		},
	}

	fakeTodosResponse := todo.TodosSerializer{
		Model: fakeTodos,
	}.Response()

	testCases := []struct {
		description    string
		argUserID      string
		stubTodos      []todo.Todo
		stubErr        error
		expectedStatus int
		expectedJSON   string
	}{
		{
			description:    "ShouldFetchTodos",
			argUserID:      "4",
			stubTodos:      fakeTodos,
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				fakeTodosResponse,
			),
		},
		{
			description:    "ShouldBeBadRequestWhenInvalidUserID",
			argUserID:      "INVALID_USER_ID",
			expectedStatus: http.StatusBadRequest,
			expectedJSON: testutil.JSONStringFromInterface(
				suite.T(),
				api.NewErrRes(user.ErrInvalidUserID),
			),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			argUserID, err := strconv.ParseInt(tc.argUserID, 10, 64)
			if err == nil {
				suite.todoRepo.
					On("AllTodosByUserID", argUserID).
					Return(tc.stubTodos, tc.stubErr)
			}

			actual := testutil.Response(
				suite.T(),
				suite.router,
				"GET",
				"api/todos?user_id="+tc.argUserID,
				nil,
			)

			suite.Equal(tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.JSONStringFromResBody(suite.T(), actual.Body)

			suite.Equal(tc.expectedJSON, actualJSON)
		})
	}
}
