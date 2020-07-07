package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/stretchr/testify/mock"
)

type TodoService struct {
	mock.Mock
}

// NewTodoService return new fake todo service
func NewTodoService() *TodoService {
	return &TodoService{}
}

func (s *TodoService) AddTodo(req todo.AddTodoRequest) (todo.TodoResponse, error) {
	args := s.Called(req)
	return args.Get(0).(todo.TodoResponse), args.Error(1)
}

func (s *TodoService) GetTodosByUserID(userID int64) ([]todo.TodoResponse, error) {
	args := s.Called(userID)
	return args.Get(0).([]todo.TodoResponse), args.Error(1)
}

func (s *TodoService) RemoveTodo(todoID string) error {
	args := s.Called(todoID)
	return args.Error(0)
}
