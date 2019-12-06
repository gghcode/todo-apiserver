package app

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
)

type appController struct {
	appService app.UsecaseInteractor
}

// NewController return new app controller
func NewController(appService app.UsecaseInteractor) api.GinController {
	return &appController{
		appService: appService,
	}
}

func (c *appController) RegisterRoutes(router gin.IRouter) {
	router.GET("/api/healthy", c.Healthy)
	router.GET("/api/version", c.AppVersion)
}

// Healthy godoc
// @Description Get server healthy
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/healthy [get]
func (c *appController) Healthy(ctx *gin.Context) {}

// Version godoc
// @Description Get server version
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/version [get]
func (c *appController) AppVersion(ctx *gin.Context) {
	ctx.String(http.StatusOK, c.appService.Version())
}
