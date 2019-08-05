package api

import "github.com/gin-gonic/gin"

// Controller is interface about api Controller.
type Controller interface {
	RegisterRoutes(router gin.IRouter)
}

// ControllerToken godoc
var ControllerToken = new(Controller)
