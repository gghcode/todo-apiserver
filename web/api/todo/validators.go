package todo

import (
	"encoding/json"
	"errors"

	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type (
	addTodoRequestDTO struct {
		Title    string `json:"title"`
		Contents string `json:"contents"`
	}
)

func validateAddTodoRequestDTO(ctx *gin.Context, req *todo.AddTodoRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto addTodoRequestDTO
	if err := ctx.ShouldBindJSON(&reqDto); errors.As(err, &jsonError) {
		return api.NewUnmarshalError(jsonError.Field, jsonError.Type.String())
	} else if err != nil {
		return err
	}

	*req = todo.AddTodoRequest{
		Title:      reqDto.Title,
		Contents:   reqDto.Contents,
		AssignorID: middleware.AuthUserID(ctx),
	}

	return validation.ValidateStruct(
		&reqDto,
		validation.Field(&reqDto.Title, validation.Required),
		validation.Field(&reqDto.Contents, validation.Required),
	)
}
