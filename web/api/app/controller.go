package app

import (
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gin-gonic/gin"
)

// Controller is app controller
type Controller struct {
	appService   app.UsecaseInteractor
	postgresConn db.GormConnection
	redisConn    db.RedisConnection
}

// NewController return new app controller
func NewController(
	appService app.UsecaseInteractor,
	postgresConn db.GormConnection,
	redisConn db.RedisConnection,
) *Controller {
	return &Controller{
		appService:   appService,
		postgresConn: postgresConn,
		redisConn:    redisConn,
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
func (c *Controller) Healthy(ctx *gin.Context) {
	healthy := true
	if !c.postgresConn.Healthy() {
		healthy = false
		// logging
	}

	if !c.redisConn.Healthy() {
		healthy = false
		// logging
	}

	if !healthy {
		ctx.Status(http.StatusServiceUnavailable)
	}
}

// AppVersion godoc
// @Description Get server version
// @Success 200 {string} string OK
// @Tags App API
// @Router /api/version [get]
func (c *Controller) AppVersion(ctx *gin.Context) {
	ctx.String(http.StatusOK, c.appService.Version())
}
