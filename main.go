package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/service"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	_ "gitlab.com/gyuhwan/apas-todo-apiserver/docs"
)

const (
	envPrefix = "REST"
)

// @title APAS TODO API
// @version 1.0
// @description This is a apas todo api server.
// @BasePath /api
func main() {
	cfg, err := config.NewBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()

	container, err := service.NewContainer(cfg)
	if err != nil {
		panic(err)
	}

	var controllers []api.Controller
	if err := container.Extract(&controllers); err != nil {
		panic(err)
	}

	router := gin.New()
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiRouter := router.Group(cfg.BasePath)
	for _, controller := range controllers {
		controller.RegisterRoutes(apiRouter)
	}

	if err := router.Run(cfg.Addr); err != nil {
		panic(err)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
