package todo_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
)

func TestTodoService_RemoveTodo(t *testing.T) {
	testCases := []struct {
		description   string
		argTodoID     string
		stubRemoveErr error
		expectedErr   error
	}{
		{
			description:   "ShouldRemoveTodo",
			argTodoID:     "abcd",
			stubRemoveErr: nil,
			expectedErr:   nil,
		},
		{
			description:   "ShouldReturnErrStub",
			argTodoID:     "",
			stubRemoveErr: fake.ErrStub,
			expectedErr:   fake.ErrStub,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakeTodoRepo := fake.NewTodoRepository()
			fakeTodoRepo.
				On("RemoveTodo", tc.argTodoID).
				Return(tc.stubRemoveErr)

			srv := todo.NewTodoService(fakeTodoRepo)
			actualErr := srv.RemoveTodo(tc.argTodoID)

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
			description:    "ShouldReturnErrStub",
			argAddTodoReq:  todo.AddTodoRequest{},
			stubAddTodoRes: todo.Todo{},
			stubAddTodoErr: fake.ErrStub,
			expectedRes:    todo.TodoResponse{},
			expectedErr:    fake.ErrStub,
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
			description:  "ShouldReturnErrStub",
			argUserID:    -1,
			stubTodosRes: nil,
			stubErr:      fake.ErrStub,
			expectedRes:  nil,
			expectedErr:  fake.ErrStub,
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
