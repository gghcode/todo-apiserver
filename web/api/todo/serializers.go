package todo

import "github.com/gghcode/apas-todo-apiserver/domain/todo"

type (
	todoResponseDTO struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		Contents   string `json:"contents"`
		AssignorID int64  `json:"assignor_id"`
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

func (s *todoResponseSerializer) Response() todoResponseDTO {
	return todoResponseDTO{
		ID:         s.model.ID,
		Title:      s.model.Title,
		Contents:   s.model.Contents,
		AssignorID: s.model.AssignorID,
	}
}
