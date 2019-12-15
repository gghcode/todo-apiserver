package app

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gin-gonic/gin"
)

// Controller is app controller
type Controller struct {
	appService app.UsecaseInteractor
}

// NewController return new app controller
func NewController(appService app.UsecaseInteractor) *Controller {
	return &Controller{
		appService: appService,
	}
}

// RegisterRoutes register handler routes.
func (c *Controller) RegisterRoutes(router gin.IRouter) {
	router.GET("/api/healthy", c.Healthy)
	router.GET("/api/version", c.AppVersion)
}

// Healthy godoc
// @Description Get server healthy
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/healthy [get]
func (c *Controller) Healthy(ctx *gin.Context) {}

// AppVersion godoc
// @Description Get server version
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/version [get]
func (c *Controller) AppVersion(ctx *gin.Context) {
	ctx.String(http.StatusOK, c.appService.Version())
}
