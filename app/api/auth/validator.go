package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

// LoginRequest godoc
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
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
