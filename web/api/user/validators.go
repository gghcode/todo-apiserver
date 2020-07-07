package user

import (
	"encoding/json"
	"errors"

	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type (
	createUserRequestDTO struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
)

func validateCreateUserRequestDTO(ctx *gin.Context, req *user.CreateUserRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto createUserRequestDTO
	if err := ctx.ShouldBindJSON(&reqDto); errors.As(err, &jsonError) {
		return api.NewUnmarshalError(jsonError.Field, jsonError.Type.String())
	} else if err != nil {
		return err
	}

	*req = user.CreateUserRequest{
		UserName: reqDto.UserName,
		Password: reqDto.Password,
	}

	return validation.ValidateStruct(
		&reqDto,
		validation.Field(&reqDto.UserName, validation.Required),
		validation.Field(&reqDto.Password, validation.Required),
	)
}
