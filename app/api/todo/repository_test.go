package todo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/todo"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
	"gitlab.com/gyuhwan/apas-todo-apiserver/internal/testutil"
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
