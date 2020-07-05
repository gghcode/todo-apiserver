package model

import (
	"time"

	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	uuid "github.com/satori/go.uuid"
)

// Todo godoc
type Todo struct {
	ID         uuid.UUID `gorm:"type:uuid"`
	Title      string    `gorm:"not null;"`
	Contents   string    `gorm:"not null;"`
	AssignorID int64     `gorm:"not null;"`
	CreatedAt  time.Time `sql:"DEFAULT:current_timestamp"`
	UpdatedAt  time.Time `sql:"DEFAULT:current_timestamp"`
}

// FromTodoEntity create todo data model from todo entity model
func FromTodoEntity(todo todo.Todo) Todo {
	return Todo{
		ID:         uuid.FromStringOrNil(todo.ID),
		Title:      todo.Title,
		Contents:   todo.Contents,
		AssignorID: todo.AssignorID,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
	}
}

// ToTodoEntity create todo entity model from todo data model
func ToTodoEntity(t Todo) todo.Todo {
	return todo.Todo{
		ID:         t.ID.String(),
		Title:      t.Title,
		Contents:   t.Contents,
		AssignorID: t.AssignorID,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}
}

// ToTodoEntityArray create todo entity array from todo data model array
func ToTodoEntityArray(todos []Todo) []todo.Todo {
	var res []todo.Todo
	for _, t := range todos {
		res = append(res, ToTodoEntity(t))
	}

	return res
}
