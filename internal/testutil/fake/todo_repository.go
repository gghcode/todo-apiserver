package fake

import (
	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/stretchr/testify/mock"
)

// TodoRepository godoc
type TodoRepository struct {
	mock.Mock
}

// AddTodo godoc
func (repo *TodoRepository) AddTodo(t todo.Todo) (todo.Todo, error) {
	args := repo.Called(t)
	return args.Get(0).(todo.Todo), args.Error(1)
}

// AllTodosByUserID godoc
func (repo *TodoRepository) AllTodosByUserID(userID int64) ([]todo.Todo, error) {
	args := repo.Called(userID)
	return args.Get(0).([]todo.Todo), args.Error(1)
}

// RemoveTodo godoc
func (repo *TodoRepository) RemoveTodo(todoID string) error {
	args := repo.Called(todoID)
	return args.Error(0)
}
