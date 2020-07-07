package repo_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/model"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/repo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type todoRepositoryIntegrationTestSuite struct {
	suite.Suite

	repo      todo.Repository
	testTodos []model.Todo
}

func TestTodoRepositoryIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(todoRepositoryIntegrationTestSuite))
}

func (suite *todoRepositoryIntegrationTestSuite) SetupTest() {
	cfg, err := config.FromEnvs()
	suite.NoError(err)

	dbConn, _, err := gorm.NewPostgresConnection(cfg)
	suite.NoError(err)

	testutil.SetupDBSandbox(suite.T(), dbConn.DB())

	suite.repo = repo.NewTodoRepository(dbConn)
	suite.testTodos = []model.Todo{
		{ID: uuid.NewV4(), Title: "test title 1", Contents: "test contents 2", AssignorID: 4},
		{ID: uuid.NewV4(), Title: "test title 2", Contents: "test contents 3", AssignorID: 4},
	}

	for i := range suite.testTodos {
		suite.NoError(dbConn.DB().Create(&suite.testTodos[i]).Error)
		suite.NoError(
			dbConn.DB().
				Where("id=?", suite.testTodos[i].ID.String()).
				First(&suite.testTodos[i]).
				Error,
		)
	}
}

func (suite *todoRepositoryIntegrationTestSuite) TestAddTodo() {
	testCases := []struct {
		description  string
		argTodo      todo.Todo
		expectedTodo func(actual todo.Todo) todo.Todo
		expectedErr  error
	}{
		{
			description: "ShouldAddTodo",
			argTodo: todo.Todo{
				ID:         uuid.NewV4().String(),
				Title:      "test title",
				Contents:   "test contents",
				AssignorID: 1,
			},
			expectedTodo: func(actual todo.Todo) todo.Todo {
				return todo.Todo{
					ID:         actual.ID,
					Title:      "test title",
					Contents:   "test contents",
					AssignorID: 1,
					CreatedAt:  actual.CreatedAt,
					UpdatedAt:  actual.UpdatedAt,
				}
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.AddTodo(tc.argTodo)

			suite.Equal(tc.expectedTodo(actualTodo), actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *todoRepositoryIntegrationTestSuite) TestAllTodosByUserID() {
	testCases := []struct {
		description   string
		argUserID     int64
		expectedTodos []todo.Todo
		expectedErr   error
	}{
		{
			description:   "ShouldReturnTodo",
			argUserID:     4,
			expectedTodos: model.ToTodoEntityArray(suite.testTodos),
			expectedErr:   nil,
		},
		{
			description:   "ShouldFetchEmpty",
			argUserID:     5,
			expectedTodos: nil,
			expectedErr:   nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodos, actualErr := suite.repo.AllTodosByUserID(tc.argUserID)

			suite.Equal(tc.expectedTodos, actualTodos)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *todoRepositoryIntegrationTestSuite) TestTodoByTodoID() {
	testCases := []struct {
		description  string
		argTodoID    string
		expectedTodo todo.Todo
		expectedErr  error
	}{
		{
			description:  "ShouldReturnTodo",
			argTodoID:    suite.testTodos[0].ID.String(),
			expectedTodo: model.ToTodoEntity(suite.testTodos[0]),
			expectedErr:  nil,
		},
		{
			description:  "ShouldReturnErrNotFoundTodo",
			argTodoID:    uuid.Nil.String(),
			expectedTodo: todo.Todo{},
			expectedErr:  todo.ErrNotFoundTodo,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actualTodo, actualErr := suite.repo.TodoByTodoID(tc.argTodoID)

			suite.Equal(tc.expectedTodo, actualTodo)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *todoRepositoryIntegrationTestSuite) TestUpdateTodo() {
	testCases := []struct {
		description string
		argTodoID   string
		argTodo     map[string]interface{}
		expected    func(args map[string]interface{}) todo.Todo
		expectedErr error
	}{
		{
			description: "ShouldUpdateTodo",
			argTodoID:   suite.testTodos[0].ID.String(),
			argTodo: map[string]interface{}{
				"Title":    "update title",
				"Contents": "update contents",
			},
			expected: func(args map[string]interface{}) todo.Todo {
				return todo.Todo{
					Title:      args["Title"].(string),
					Contents:   args["Contents"].(string),
					ID:         suite.testTodos[0].ID.String(),
					AssignorID: suite.testTodos[0].AssignorID,
					CreatedAt:  suite.testTodos[0].CreatedAt,
					UpdatedAt:  suite.testTodos[0].UpdatedAt,
				}
			},
			expectedErr: nil,
		},
		{
			description: "ShouldUpdateTodoWithBlackField",
			argTodoID:   suite.testTodos[1].ID.String(),
			argTodo: map[string]interface{}{
				"Title": "",
			},
			expected: func(args map[string]interface{}) todo.Todo {
				return todo.Todo{
					Title:      args["Title"].(string),
					Contents:   suite.testTodos[1].Contents,
					ID:         suite.testTodos[1].ID.String(),
					AssignorID: suite.testTodos[1].AssignorID,
					CreatedAt:  suite.testTodos[1].CreatedAt,
					UpdatedAt:  suite.testTodos[1].UpdatedAt,
				}
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotFoundTodo",
			argTodoID:   uuid.Nil.String(),
			argTodo:     map[string]interface{}{},
			expected: func(args map[string]interface{}) todo.Todo {
				return todo.Todo{}
			},
			expectedErr: todo.ErrNotFoundTodo,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.repo.UpdateTodo(
				tc.argTodoID,
				tc.argTodo,
			)

			suite.Equal(tc.expected(tc.argTodo), actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *todoRepositoryIntegrationTestSuite) TestRemoveTodo() {
	testCases := []struct {
		description string
		argTodoID   uuid.UUID
		expectedErr error
	}{
		{
			description: "ShouldRemoveTodo",
			argTodoID:   suite.testTodos[0].ID,
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotFoundTodo",
			argTodoID:   uuid.Nil,
			expectedErr: todo.ErrNotFoundTodo,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			argTodoID := tc.argTodoID.String()
			actualErr := suite.repo.RemoveTodo(argTodoID)

			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
