package user

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gin-gonic/gin"
)

// CreateUserRequest is dto that contains info that require to create user.
type CreateUserRequest struct {
	UserName string `json:"username" validate:"required,min=4,max=100"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// CreateUserRequestValidator godoc
type CreateUserRequestValidator struct {
	Model CreateUserRequest
}

// Bind godoc
func (v *CreateUserRequestValidator) Bind(ctx *gin.Context) error {
	if err := ctx.ShouldBindJSON(&v.Model); err != nil {
		return api.NewHandledError(http.StatusBadRequest, err)
	}

	return api.Validate(v.Model)
}
