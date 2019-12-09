package auth

import (
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Controller godoc
type Controller struct {
	service auth.UsecaseInteractor
}

// NewController return new auth controller instance.
func NewController(service auth.UsecaseInteractor) *Controller {
	return &Controller{
		service: service,
	}
}

// RegisterRoutes register handler routes.
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.POST("api/auth/token", controller.issueToken)
	router.POST("api/auth/refresh", controller.refreshToken)
}

// @Description Get new access token by refreshtoken
// @Accept json
// @Produce json
// @Param payload body auth.AccessTokenByRefreshRequest true "payload"
// @Success 200 {object} auth.tokenResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid payload"
// @Failure 401 {object} api.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /api/auth/refresh [post]
func (controller *Controller) refreshToken(ctx *gin.Context) {
	var req auth.AccessTokenByRefreshRequest
	if err := validateAccessTokenByRefreshRequest(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponse(err))
		return
	}

	res, err := controller.service.RefreshToken(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponse(err))
		return
	}

	serializer := newTokenResponseSerializer(res)

	ctx.JSON(http.StatusOK, serializer.Response())
}

// @Description Issue new token
// @Accept json
// @Produce json
// @Param payload body auth.LoginRequest true "payload"
// @Success 200 {object} auth.TokenResponse "ok"
// @Failure 400 {object} api.ErrorResponse "Invalid payload"
// @Failure 401 {object} api.ErrorResponse "Invalid credential"
// @Tags Auth API
// @Router /api/auth/token [post]
func (controller *Controller) issueToken(ctx *gin.Context) {
	var req auth.LoginRequest
	if err := validateLoginRequestDto(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponse(err))
		return
	}

	res, err := controller.service.IssueToken(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponse(err))
		return
	}

	serializer := newTokenResponseSerializer(res)

	ctx.JSON(http.StatusOK, serializer.Response())
	// loginRequestValidator := NewLoginRequestValidator()
	// if err := loginRequestValidator.Bind(ctx); err != nil {
	// 	api.WriteErrorResponse(ctx, err)
	// 	return
	// }

	// token, err := controller.service.IssueToken(loginRequestValidator.Model)
	// if err != nil {
	// 	api.WriteErrorResponse(ctx, err)
	// 	return
	// }

	// ctx.JSON(http.StatusOK, token)
}
