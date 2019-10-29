package main

import (
	"github.com/gghcode/apas-todo-apiserver/app"
	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/auth"
	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gghcode/apas-todo-apiserver/config"
	_ "github.com/gghcode/apas-todo-apiserver/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	envPrefix = "REST"
)

// @title APAS TODO API
// @version 1.0
// @description This is a apas todo api server.
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewViperBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()

	container, err := app.NewContainer(cfg)
	if err != nil {
		panic(err)
	}

	var controllers []api.Controller
	if err := container.Extract(&controllers); err != nil {
		panic(err)
	}

	router := gin.New()
	router.Use(middleware.NewCors(cfg.Cors))
	router.Use(middleware.AddJwtAuthHandler(cfg.Jwt, &auth.JwtAuthHandler))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, controller := range controllers {
		controller.RegisterRoutes(router.Group(""))
	}

	if err := router.Run(cfg.Addr); err != nil {
		panic(err)
	}
}
