package auth

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/web/api"

	"github.com/gin-gonic/gin"
)

// Controller godoc
type Controller struct {
	service auth.UseCase
}

// NewController return new auth controller instance.
func NewController(service auth.UseCase) *Controller {
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
// @Param payload body auth.accessTokenByRefreshRequestDTO true "payload"
// @Success 200 {object} auth.tokenResponseDTO "ok"
// @Failure 400 {object} api.ErrorResponseDTO "Invalid payload"
// @Failure 401 {object} api.ErrorResponseDTO "Invalid credential"
// @Tags Auth API
// @Router /api/auth/refresh [post]
func (controller *Controller) refreshToken(ctx *gin.Context) {
	var req auth.AccessTokenByRefreshRequest
	if err := validateAccessTokenByRefreshRequestDTO(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponseDTO(err))
		return
	}

	res, err := controller.service.RefreshToken(req)
	if err == auth.ErrNotStoredToken {
		ctx.JSON(http.StatusUnauthorized, api.MakeErrorResponseDTO(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	serializer := newTokenResponseSerializer(res)

	ctx.JSON(http.StatusOK, serializer.Response())
}

// @Description Issue new token
// @Accept json
// @Produce json
// @Param payload body auth.loginRequestDTO true "payload"
// @Success 200 {object} auth.tokenResponseDTO "ok"
// @Failure 400 {object} api.ErrorResponseDTO "Invalid payload"
// @Failure 401 {object} api.ErrorResponseDTO "Invalid credential"
// @Tags Auth API
// @Router /api/auth/token [post]
func (controller *Controller) issueToken(ctx *gin.Context) {
	var req auth.LoginRequest
	if err := validateLoginRequestDTO(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, api.MakeErrorResponseDTO(err))
		return
	}

	res, err := controller.service.IssueToken(req)
	if err == auth.ErrInvalidCredential {
		ctx.JSON(http.StatusUnauthorized, api.MakeErrorResponseDTO(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, api.MakeErrorResponseDTO(err))
		return
	}

	serializer := newTokenResponseSerializer(res)

	ctx.JSON(http.StatusOK, serializer.Response())
}
