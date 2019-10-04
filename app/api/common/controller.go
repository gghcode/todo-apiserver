package common

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/app/loader"
	"github.com/gin-gonic/gin"
)

// Controller godoc
type Controller struct {
	appVersion string
}

// NewController godoc
func NewController(appVersionLoader loader.VersionLoader) *Controller {
	return &Controller{
		appVersion: appVersionLoader.GetVersion(),
	}
}

// RegisterRoutes godocs
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.GET("api/healthy", controller.Healthy)
	router.GET("api/version", controller.Version)
}

// Version godoc
// @Description Get server version
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/version [get]
func (controller *Controller) Version(ctx *gin.Context) {
	ctx.String(http.StatusOK, controller.appVersion)
}

// Healthy godoc
// @Description Get server healthy
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/healthy [get]
func (controller *Controller) Healthy(ctx *gin.Context) {}
