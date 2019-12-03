package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/stretchr/testify/mock"
)

// Todo2Repository godoc
type Todo2Repository struct {
	mock.Mock
}

// NewTodoRepository return fake todo repository
func NewTodoRepository() *Todo2Repository {
	return &Todo2Repository{}
}

// AddTodo godoc
func (repo *Todo2Repository) AddTodo(t todo.Todo) (todo.Todo, error) {
	args := repo.Called(t)
	return args.Get(0).(todo.Todo), args.Error(1)
}

// AllTodosByUserID godoc
func (repo *Todo2Repository) AllTodosByUserID(userID int64) ([]todo.Todo, error) {
	args := repo.Called(userID)
	return args.Get(0).([]todo.Todo), args.Error(1)
}

// TodoByTodoID godoc
func (repo *Todo2Repository) TodoByTodoID(todoID string, todo *todo.Todo) error {
	args := repo.Called(todoID, todo)
	return args.Error(0)
}

// UpdateTodo godoc
func (repo *Todo2Repository) UpdateTodo(todoID string, todoData map[string]interface{}) (todo.Todo, error) {
	args := repo.Called(todoID, todoData)
	return args.Get(0).(todo.Todo), args.Error(1)
}

// RemoveTodo godoc
func (repo *Todo2Repository) RemoveTodo(todoID string) error {
	args := repo.Called(todoID)
	return args.Error(0)
}
