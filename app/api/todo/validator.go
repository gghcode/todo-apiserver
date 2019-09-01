package todo

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

// AddTodoRequest godoc
// struct tag use swagger
type AddTodoRequest struct {
	Title    string `json:"title" validate:"required,min=1"`
	Contents string `json:"contents" validate:"required,min=1"`
}

// AddTodoValidator godoc
type AddTodoValidator struct {
	Model AddTodoRequest
}

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
