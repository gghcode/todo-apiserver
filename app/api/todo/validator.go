package todo

import (
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
		Title    string    `json:"title" validate:"min=1"`
		Contents string    `json:"contents" validate:"min=1"`
		DueDate  time.Time `json:"due_date"`
	}
)
