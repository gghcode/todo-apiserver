package api

import "github.com/gin-gonic/gin"

// GinController is interface about api Controller.
type GinController interface {
	RegisterRoutes(gin.IRouter)
}
