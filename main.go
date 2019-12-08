package main

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/jwt"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/repository"
	infraSecurity "github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/gghcode/apas-todo-apiserver/web"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	webApp "github.com/gghcode/apas-todo-apiserver/web/api/app"
	webAuth "github.com/gghcode/apas-todo-apiserver/web/api/auth"

	"github.com/defval/inject/v2"
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

	ginRouter.Run(cfg.Addr)
}

func setupContainer(cfg config.Configuration) *inject.Container {
	container := inject.New(
		inject.Provide(func() config.Configuration {
			return cfg
		}),

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
		inject.Provide(web.NewMux),
	)

	return container
}
