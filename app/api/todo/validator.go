package todo

import (
	"database/sql"
	"time"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gin-gonic/gin"
)

type (
	// AddTodoRequest godoc
	// struct tag use swagger
	AddTodoRequest struct {
		Title    string `json:"title" validate:"required,min=1"`
		Contents string `json:"contents" validate:"required,min=1"`
	}

	// AddTodoValidator godoc
	AddTodoValidator struct {
		Model AddTodoRequest
	}
)

// NewAddTodoValidator godoc
func NewAddTodoValidator() *AddTodoValidator {
	return &AddTodoValidator{}
}

// Bind godoc
func (v *AddTodoValidator) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&v.Model); err != nil {
		return err
	}

	return api.Validate(v.Model)
}

type (
	// UpdateTodoRequest godoc
	UpdateTodoRequest struct {
		Title    string    `json:"title"`
		Contents string    `json:"contents"`
		DueDate  time.Time `json:"due_date"`
	}

	// UpdateTodoValidator godoc
	UpdateTodoValidator struct {
		Model UpdateTodoRequest
	}
)

// NewUpdateTodoValidator godoc
func NewUpdateTodoValidator() *UpdateTodoValidator {
	return &UpdateTodoValidator{}
}

// Bind godoc
func (v *UpdateTodoValidator) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&v.Model); err != nil {
		return err
	}

	return api.Validate(v.Model)
}

// Entity convert to entity from request
func (model UpdateTodoRequest) Entity() Todo {
	return Todo{
		Title:    model.Title,
		Contents: model.Contents,
		DueDate: sql.NullTime{
			Time: model.DueDate,
		},
	}
}
