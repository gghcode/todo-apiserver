package todo_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
)

func TestTodoService_GetTodosByUserID(t *testing.T) {
	testCases := []struct {
		description  string
		argUserID    int64
		stubTodosRes []todo.Todo
		stubErr      error
		expectedRes  []todo.TodoResponse
		expectedErr  error
	}{
		{
			description: "ShouldFetchTodos",
			argUserID:   1,
			stubTodosRes: []todo.Todo{
				{
					ID:       uuid.Nil,
					Title:    "test title1",
					Contents: "test contents1",
				},
				{
					ID:       uuid.Nil,
					Title:    "test title2",
					Contents: "test contents2",
				},
			},
			stubErr: nil,
			expectedRes: []todo.TodoResponse{
				{
					ID:       uuid.Nil.String(),
					Title:    "test title1",
					Contents: "test contents1",
				},
				{
					ID:       uuid.Nil.String(),
					Title:    "test title2",
					Contents: "test contents2",
				},
			},
			expectedErr: nil,
		},
		{
			description:  "ShouldReturnErrFake",
			argUserID:    -1,
			stubTodosRes: nil,
			stubErr:      fake.ErrFake,
			expectedRes:  nil,
			expectedErr:  fake.ErrFake,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakeTodoRepo := fake.NewTodoRepository()
			fakeTodoRepo.
				On("AllTodosByUserID", tc.argUserID).
				Return(tc.stubTodosRes, tc.stubErr)

			srv := todo.NewTodoService(fakeTodoRepo)

			actualRes, actualErr := srv.GetTodosByUserID(tc.argUserID)

			assert.Equal(t, tc.expectedRes, actualRes)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func TestTodoService_AddTodo(t *testing.T) {
	testCases := []struct {
		description    string
		argAddTodoReq  todo.AddTodoRequest
		stubAddTodoRes todo.Todo
		stubAddTodoErr error
		expectedRes    todo.TodoResponse
		expectedErr    error
	}{
		{
			description: "ShouldAddTodo",
			argAddTodoReq: todo.AddTodoRequest{
				Title:    "test title",
				Contents: "test contents",
			},
			stubAddTodoRes: todo.Todo{
				Title:    "test title",
				Contents: "test contents",
			},
			stubAddTodoErr: nil,
			expectedRes: todo.TodoResponse{
				Title:    "test title",
				Contents: "test contents",
			},
			expectedErr: nil,
		},
		{
			description:    "ShouldReturnErrFake",
			argAddTodoReq:  todo.AddTodoRequest{},
			stubAddTodoRes: todo.Todo{},
			stubAddTodoErr: fake.ErrFake,
			expectedRes:    todo.TodoResponse{},
			expectedErr:    fake.ErrFake,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakeTodoRepo := fake.NewTodoRepository()
			fakeTodoRepo.
				On("AddTodo", mock.Anything).
				Return(tc.stubAddTodoRes, tc.stubAddTodoErr)

			todoService := todo.NewTodoService(fakeTodoRepo)

			actualRes, actualErr := todoService.AddTodo(tc.argAddTodoReq)

			assert.Equal(t, tc.expectedErr, actualErr)

			assert.Equal(t, tc.expectedRes.Title, actualRes.Title)
			assert.Equal(t, tc.expectedRes.Contents, actualRes.Contents)
			assert.Equal(t, tc.expectedRes.AssignorID, actualRes.AssignorID)
			assert.Equal(t, tc.expectedRes.DueDate, actualRes.DueDate)
		})
	}
}
