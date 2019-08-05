package common

import "github.com/gin-gonic/gin"

// Controller godoc
type Controller struct {
}

// NewController godoc
func NewController() *Controller {
	return &Controller{}
}

// RegisterRoutes godocs
func (controller *Controller) RegisterRoutes(router gin.IRouter) {
	router.GET("/healthy", controller.Healthy)
}

// Healthy godoc
// @Description Get server healthy
// @Success 200 {string} string OK
// @Tags App API
// @Router /healthy [get]
func (controller *Controller) Healthy(ctx *gin.Context) {}
