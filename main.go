package main

import (
	"github.com/gghcode/apas-todo-apiserver/config"
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
	config.NewViperBuilder().
		AddConfigFile("config.yaml", true).
		BindEnvs(envPrefix).
		Build()
}

// func setupContainer(cfg config.Configuration) (*inject.Container, error) {
// 	container, err := inject.New(
// 		inject.Provide(func() config.Configuration {
// 			return cfg
// 		}),

// 		inject.Provide(db.NewPostgresConn),
// 		inject.Provide(db.NewRedisConn),
// 		inject.Provide(func() security.Passport {
// 			return security.NewBcryptPassport(12)
// 		}),
// 		inject.Provide(func() afero.Fs {
// 			return afero.NewOsFs()
// 		}),

// 		inject.Provide(loader.NewVersionLoader),

// 		inject.Provide(user.NewRepository),
// 		inject.Provide(todo.NewRepository),
// 		inject.Provide(auth.NewRedisTokenRepository),

// 		inject.Provide(auth.NewService),
// 		inject.Provide(func() auth.CreateAccessTokenHandlerFactory {
// 			return auth.CreateAccessTokenFactory
// 		}),
// 		inject.Provide(func() auth.CreateRefreshTokenHandlerFactory {
// 			return auth.CreateRefreshTokenFactory
// 		}),

// 		inject.Provide(common.NewController, inject.As(api.ControllerToken)),
// 		inject.Provide(user.NewController, inject.As(api.ControllerToken)),
// 		inject.Provide(todo.NewController, inject.As(api.ControllerToken)),
// 		inject.Provide(auth.NewController, inject.As(api.ControllerToken)),
// 	)

// 	return container, err
// }
