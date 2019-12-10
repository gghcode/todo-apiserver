package repository_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
)

type todoRepositoryIntegrationTestSuite struct {
	suite.Suite

	repo      todo.Repository
	dbCleanup func()
}

func TestTodoRepositoryIntegrationTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(todoRepositoryIntegrationTestSuite))
}

func (suite *todoRepositoryIntegrationTestSuite) SetupTest() {
	cfg, err := config.NewViperBuilder().
		BindEnvs("TEST").
		Build()
	suite.NoError(err)

	dbConn, err := db.NewPostgresConn(cfg)
	suite.NoError(err)

	suite.dbCleanup = testutil.DbCleanupFunc(dbConn.DB())
	suite.repo = repository.NewGormTodoRepository(dbConn)
}

func (suite *todoRepositoryIntegrationTestSuite) TearDownTest() {
	suite.dbCleanup()
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
				ID:         uuid.NewV4(),
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
