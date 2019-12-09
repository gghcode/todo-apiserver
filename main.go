package main

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	_ "github.com/gghcode/apas-todo-apiserver/docs"
	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	infraSecurity "github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webApp "github.com/gghcode/apas-todo-apiserver/web/api/app"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/defval/inject"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
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

	if err != nil {
		panic(err)
	}

	var ginRouter *gin.Engine

	c := setupContainer(cfg)
	if err := c.Extract(&ginRouter); err != nil {
		panic(err)
	}

	ginRouter.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	ginRouter.Run(cfg.Addr)
}

func setupContainer(cfg config.Configuration) *inject.Container {
	provideConfig := func() config.Configuration {
		return cfg
	}

	container, _ := inject.New(
		inject.Provide(provideConfig),

		inject.Provide(db.NewPostgresConn),
		inject.Provide(db.NewRedisConn),
		inject.Provide(infraSecurity.NewBcryptPassport),

		inject.Provide(afero.NewOsFs),
		inject.Provide(file.NewAferoFileReader),
		inject.Provide(app.NewService),
		inject.Provide(webApp.NewController, inject.As(new(api.GinController))),

		inject.Provide(repository.NewRedisTokenRepository),
		inject.Provide(repository.NewUserRepository),
		inject.Provide(auth.NewService),
		inject.Provide(jwt.NewJwtAccessTokenHandlerFactory),
		inject.Provide(jwt.NewJwtRefreshTokenHandlerfactory),
		inject.Provide(webAuth.NewController, inject.As(new(api.GinController))),
		inject.Provide(newGinRouter),
	)

	return container
}

func newGinRouter(controllers []api.GinController) *gin.Engine {
	router := gin.New()

	for _, c := range controllers {
		c.RegisterRoutes(router)
	}

	return router
}
