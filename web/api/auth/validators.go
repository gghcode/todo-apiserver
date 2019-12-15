package auth

import (
	"encoding/json"
	"errors"

	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/web/api"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v3"
)

type (
	loginRequestDTO struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	accessTokenByRefreshRequestDTO struct {
		Token string `json:"token"`
	}
)

func validateLoginRequestDTO(ctx *gin.Context, req *auth.LoginRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto loginRequestDTO
	if err := ctx.ShouldBindJSON(&reqDto); errors.As(err, &jsonError) {
		return api.NewUnmarshalError(jsonError.Field, jsonError.Type.String())
	} else if err != nil {
		return err
	}

	*req = auth.LoginRequest{
		Username: reqDto.Username,
		Password: reqDto.Password,
	}

	return validation.ValidateStruct(
		&reqDto,
		validation.Field(&reqDto.Username, validation.Required),
		validation.Field(&reqDto.Password, validation.Required),
	)
}

func validateAccessTokenByRefreshRequestDTO(ctx *gin.Context, req *auth.AccessTokenByRefreshRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto accessTokenByRefreshRequestDTO
	if err := ctx.ShouldBindJSON(&reqDto); errors.As(err, &jsonError) {
		return api.NewUnmarshalError(jsonError.Field, jsonError.Type.String())
	} else if err != nil {
		return err
	}

	*req = auth.AccessTokenByRefreshRequest{
		Token: reqDto.Token,
	}

	return validation.ValidateStruct(
		&reqDto,
		validation.Field(&reqDto.Token, validation.Required),
	)
}
