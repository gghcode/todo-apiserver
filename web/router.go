package web

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/web/api"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// NewGinRouter return gin router
func NewGinRouter(
	cfg config.Configuration,
	middlwares []gin.HandlerFunc,
	controllers []api.GinController,
) *gin.Engine {
	router := gin.New()
	for _, m := range middlwares {
		router.Use(m)
	}

	for _, c := range controllers {
		c.RegisterRoutes(router.Group(""))
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
