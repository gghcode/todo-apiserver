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
	loginRequestDto struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	accessTokenByRefreshRequestDto struct {
		Token string `json:"token"`
	}
)

func validateLoginRequestDto(ctx *gin.Context, req *auth.LoginRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto loginRequestDto
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

func validateAccessTokenByRefreshRequest(ctx *gin.Context, req *auth.AccessTokenByRefreshRequest) error {
	var jsonError *json.UnmarshalTypeError

	var reqDto accessTokenByRefreshRequestDto
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
