package todo_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type RepositoryIntegration struct {
	suite.Suite

	pgConn         *db.PostgresConn
	dbCleanup      func()
	todoRepository todo.Repository

	testTodos []todo.Todo
}

func TestTodoRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(RepositoryIntegration))
}

func (suite *RepositoryIntegration) SetupTest() {
	cfg, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	suite.NoError(err)

	suite.pgConn, err = db.NewPostgresConn(cfg)
	suite.pgConn.DB().LogMode(false)
	suite.dbCleanup = testutil.DbCleanupFunc(suite.pgConn.DB())
	suite.todoRepository = todo.NewRepository(suite.pgConn)

	suite.testTodos = []todo.Todo{
		{Title: "test title 1", Contents: "test contents 2", AssignorID: 4},
		{Title: "test title 2", Contents: "test contents 3", AssignorID: 4},
	}

	for i := range suite.testTodos {
		err := suite.pgConn.DB().Create(&suite.testTodos[i]).Error
		suite.NoError(err)

		err = suite.pgConn.DB().
			Where("id=?", suite.testTodos[i].ID.String()).
			First(&suite.testTodos[i]).
			Error

		suite.NoError(err)
	}
}

func (suite *RepositoryIntegration) TearDownTest() {
	suite.dbCleanup()
	suite.pgConn.Close()
}

func (suite *RepositoryIntegration) TestAddTodo() {
	fakeTodo := todo.Todo{
		Title:      "fake title",
		Contents:   "fake contents",
		AssignorID: 100,
	}

	testCases := []struct {
		description string
		argTodo     todo.Todo
		expected    todo.Todo
		expectedErr error
	}{
		{
			description: "ShouldAddTodo",
			argTodo:     fakeTodo,
			expected:    fakeTodo,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.todoRepository.AddTodo(tc.argTodo)

			suite.NotEqual(tc.expected.ID, actual.ID)
			suite.NotEqual(tc.expected.CreatedAt, actual.CreatedAt)

			suite.Equal(tc.expected.Title, actual.Title)
			suite.Equal(tc.expected.Contents, actual.Contents)
			suite.Equal(tc.expected.AssignorID, actual.AssignorID)

			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *RepositoryIntegration) TestTodoByTodoID() {
	testCases := []struct {
		description string
		argTodoID   string
		expected    todo.Todo
		expectedErr error
	}{
		{
			description: "ShouldReturnTodo",
			argTodoID:   suite.testTodos[0].ID.String(),
			expected:    suite.testTodos[0],
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotFoundTodo",
			argTodoID:   uuid.Nil.String(),
			expected:    todo.EmptyTodo,
			expectedErr: todo.ErrNotFoundTodo,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			var actual todo.Todo
			actualErr := suite.todoRepository.
				TodoByTodoID(tc.argTodoID, &actual)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *RepositoryIntegration) TestUpdateTodo() {
	expectedTodo := map[string]interface{}{
		"Title":    "update title",
		"Contents": "update contents",
		"DueDate": sql.NullTime{
			Time:  time.Unix(100000, 0),
			Valid: true,
		},
	}

	expectedTodoWithBlackField := map[string]interface{}{
		"Title": "",
	}

	testCases := []struct {
		description string
		argTodoID   string
		argTodo     map[string]interface{}
		expected    todo.Todo
		expectedErr error
	}{
		{
			description: "ShouldUpdateTodo",
			argTodoID:   suite.testTodos[0].ID.String(),
			argTodo:     expectedTodo,
			expected: todo.Todo{
				Title:      expectedTodo["Title"].(string),
				Contents:   expectedTodo["Contents"].(string),
				DueDate:    expectedTodo["DueDate"].(sql.NullTime),
				ID:         suite.testTodos[0].ID,
				AssignorID: suite.testTodos[0].AssignorID,
				CreatedAt:  suite.testTodos[0].CreatedAt,
				UpdatedAt:  suite.testTodos[0].UpdatedAt,
			},
			expectedErr: nil,
		},
		{
			description: "ShouldUpdateTodoWithBlackField",
			argTodoID:   suite.testTodos[1].ID.String(),
			argTodo:     expectedTodoWithBlackField,
			expected: todo.Todo{
				Title:      expectedTodoWithBlackField["Title"].(string),
				Contents:   suite.testTodos[1].Contents,
				DueDate:    suite.testTodos[1].DueDate,
				ID:         suite.testTodos[1].ID,
				AssignorID: suite.testTodos[1].AssignorID,
				CreatedAt:  suite.testTodos[1].CreatedAt,
				UpdatedAt:  suite.testTodos[1].UpdatedAt,
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotFoundTodo",
			argTodoID:   uuid.Nil.String(),
			argTodo:     expectedTodo,
			expected:    todo.EmptyTodo,
			expectedErr: todo.ErrNotFoundTodo,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.todoRepository.UpdateTodo(tc.argTodoID, tc.argTodo)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *RepositoryIntegration) TestRemoveTodo() {
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
			actualErr := suite.todoRepository.RemoveTodo(argTodoID)

			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}

func (suite *RepositoryIntegration) TestAllTodosByUserID() {
	testCases := []struct {
		description string
		argUserID   int64
		expected    []todo.Todo
		expectedErr error
	}{
		{
			description: "ShouldFetchTodos",
			argUserID:   4,
			expected:    suite.testTodos,
			expectedErr: nil,
		},
		{
			description: "ShouldFetchEmpty",
			argUserID:   5,
			expected:    []todo.Todo{},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.description, func() {
			actual, actualErr := suite.todoRepository.AllTodosByUserID(tc.argUserID)

			suite.Equal(tc.expected, actual)
			suite.Equal(tc.expectedErr, actualErr)
		})
	}
}
