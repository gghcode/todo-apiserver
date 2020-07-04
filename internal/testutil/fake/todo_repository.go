package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/stretchr/testify/mock"
)

// TodoRepository godoc
type TodoRepository struct {
	mock.Mock
}

// NewTodoRepository return fake todo repository
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{}
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

// TodoByTodoID godoc
func (repo *TodoRepository) TodoByTodoID(todoID string) (todo.Todo, error) {
	args := repo.Called(todoID)
	return args.Get(0).(todo.Todo), args.Error(1)
}

// UpdateTodo godoc
func (repo *TodoRepository) UpdateTodo(todoID string, todoData map[string]interface{}) (todo.Todo, error) {
	args := repo.Called(todoID, todoData)
	return args.Get(0).(todo.Todo), args.Error(1)
}

// RemoveTodo godoc
func (repo *TodoRepository) RemoveTodo(todoID string) error {
	args := repo.Called(todoID)
	return args.Error(0)
}
