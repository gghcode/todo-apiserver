package todo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/tool"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
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
		Title    *string `json:"title"`
		Contents *string `json:"contents"`
		DueDate  *string `json:"due_date"`
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
	var jsonErr *json.UnmarshalTypeError

	err := ctx.ShouldBindJSON(&v.Model)
	if errors.As(err, &jsonErr) {
		return api.JsonTypeError{
			Value: jsonErr.Value,
			Field: jsonErr.Field,
		}
	} else if err != nil {
		return err
	}

	return validation.ValidateStruct(&v.Model,
		validation.Field(&v.Model.Contents, validation.Length(2, 50)),
		validation.Field(&v.Model.Title, validation.Length(2, 100)),
	)
}

// Map godoc
func (model UpdateTodoRequest) Map() map[string]interface{} {
	result := map[string]interface{}{}

	if model.Title != nil {
		result["title"] = *model.Title
	}

	if model.Contents != nil {
		result["contents"] = *model.Contents
	}

	if model.DueDate != nil {
		result["due_date"] = *model.DueDate
	}

	return result
}

// Entity convert to entity from request
func (model UpdateTodoRequest) Entity() Todo {
	sqlNullTime := sql.NullTime{Valid: false}
	if model.DueDate != nil {
		parsedTime, err := tool.ParseTime(*model.DueDate)
		if err == nil {
			sqlNullTime.Time = parsedTime
			sqlNullTime.Valid = true
		}
	}
	fmt.Println(sqlNullTime)
	// parsedTime, err := tool.ParseTime(*model.DueDate)
	// sqlNullTime := sql.NullTime{Time: parsedTime, Valid: true}
	// if err != nil {
	// 	fmt.Println(err)
	// 	sqlNullTime.Time = time.Time{}
	// 	sqlNullTime.Valid = false
	// }
	// fmt.Println(sqlNullTime)
	return Todo{
		Title:    *model.Title,
		Contents: *model.Contents,
		DueDate:  sqlNullTime,
	}
}
