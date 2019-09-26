package todo_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestTodoSerializer(t *testing.T) {
	fakeTodoModel := todo.Todo{
		ID:       uuid.NewV4(),
		Title:    "fake title",
		Contents: "fake contents",
	}

	testCases := []struct {
		description  string
		argTodoModel todo.Todo
		expected     todo.TodoResponse
	}{
		{
			description:  "ShouldBeEqual",
			argTodoModel: fakeTodoModel,
			expected: todo.TodoResponse{
				ID:       fakeTodoModel.ID.String(),
				Title:    fakeTodoModel.Title,
				Contents: fakeTodoModel.Contents,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			todoSerializer := todo.TodoSerializer{
				Model: tc.argTodoModel,
			}

			actual := todoSerializer.Response()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestTodosSerializer(t *testing.T) {
	fakeTodoModels := []todo.Todo{
		{ID: uuid.NewV4(), Title: "fake title", Contents: "fake contents"},
		{ID: uuid.NewV4(), Title: "fake2 title", Contents: "fake2 contents"},
		{ID: uuid.NewV4(), Title: "fake3 title", Contents: "fake3 contents"},
	}

	var fakeTodosResponse []todo.TodoResponse
	for _, model := range fakeTodoModels {
		fakeTodosResponse = append(fakeTodosResponse, todo.TodoResponse{
			ID:       model.ID.String(),
			Title:    model.Title,
			Contents: model.Contents,
		})
	}

	testCases := []struct {
		description   string
		argTodoModels []todo.Todo
		expected      []todo.TodoResponse
	}{
		{
			description:   "ShouldBeEqual",
			argTodoModels: fakeTodoModels,
			expected:      fakeTodosResponse,
		},
		{
			description:   "ShouldBeEqualWhenEmptySlice",
			argTodoModels: []todo.Todo{},
			expected:      []todo.TodoResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			todosSerializer := todo.TodosSerializer{
				Model: tc.argTodoModels,
			}

			actual := todosSerializer.Response()

			assert.Equal(t, tc.expected, actual)
		})
	}
}
