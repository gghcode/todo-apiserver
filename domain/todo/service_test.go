package todo_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddTodo(t *testing.T) {
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
