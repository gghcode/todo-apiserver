package todo

import "github.com/gghcode/apas-todo-apiserver/domain/todo"

type (
	todoResponse struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Contents string `json:"contents"`
	}

	todoResponseSerializer struct {
		model todo.TodoResponse
	}
)

func newTodoResponseSerializer(model todo.TodoResponse) *todoResponseSerializer {
	return &todoResponseSerializer{
		model: model,
	}
}

func (s *todoResponseSerializer) Response() todoResponse {
	return todoResponse{
		ID:       s.model.ID,
		Title:    s.model.Title,
		Contents: s.model.Contents,
	}
}
