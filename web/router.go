package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
) (*http.Server, func()) {
	router := gin.New()
	for _, m := range middlwares {
		router.Use(m)
	}

	for _, c := range controllers {
		c.RegisterRoutes(router.Group(""))
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	gracefullyShutdownFunc := func() {
		fmt.Println("Shutdown server...")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.GracefulShutdownTimeoutSec)*time.Second,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("Server shutdown: ", err)
		}

		fmt.Println("Shutdown was successful")
	}

	return srv, gracefullyShutdownFunc
}
