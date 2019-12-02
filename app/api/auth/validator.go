package auth

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gin-gonic/gin"
)

// RefreshTokenRequestValidator godoc
type RefreshTokenRequestValidator struct {
	Model AccessTokenByRefreshRequest
}

// Bind godoc
func (v *RefreshTokenRequestValidator) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&v.Model); err != nil {
		return api.NewHandledError(http.StatusBadRequest, err)
	}

	return api.Validate(v.Model)
}

// LoginRequestValidator godoc
type LoginRequestValidator struct {
	Model LoginRequest
}

// NewLoginRequestValidator godoc
func NewLoginRequestValidator() *LoginRequestValidator {
	return &LoginRequestValidator{
		Model: LoginRequest{},
	}
}

// Bind godoc
func (v *LoginRequestValidator) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&v.Model); err != nil {
		return api.NewHandledError(http.StatusBadRequest, err)
	}

	return api.Validate(v.Model)
}
