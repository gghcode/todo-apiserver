package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

// Controller godoc
type Controller struct {
	service Service
}

// NewController return new auth controller instance.
func NewController(service Service) *Controller {
	return &Controller{
		service: service,
	}
}

// APIPath is path prefix
const APIPath = "auth"

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.Handle("POST", APIPath+"/token", controller.issueToken)
	router.Handle("POST", APIPath+"/refresh", controller.refreshToken)
}

// @Description Get new access token by refreshtoken
// @Accept json
// @Produce json
// @Param payload body auth.AccessTokenByRefreshRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid payload"
// @Failure 401 {object} common.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /auth/refresh [post]
func (controller *Controller) refreshToken(ctx *gin.Context) {
	refreshTokenRequestValidator := RefreshTokenRequestValidator{}
	if err := refreshTokenRequestValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	token, err := controller.service.RefreshToken(refreshTokenRequestValidator.Model)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}

// @Description Issue new token
// @Accept json
// @Produce json
// @Param payload body auth.LoginRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} common.ErrorResponse "Invalid payload"
// @Failure 401 {object} common.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /auth/token [post]
func (controller *Controller) issueToken(ctx *gin.Context) {
	loginRequestValidator := NewLoginRequestValidator()
	if err := loginRequestValidator.Bind(ctx); err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	token, err := controller.service.IssueToken(loginRequestValidator.Model)
	if err != nil {
		api.WriteErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, token)
}
