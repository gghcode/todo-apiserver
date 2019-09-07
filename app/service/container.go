package service

import (
	"github.com/defval/inject"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/auth"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/todo"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/infra"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
)

// NewContainer godoc
func NewContainer(cfg config.Configuration) (*inject.Container, error) {
	container, err := inject.New(
		inject.Provide(func() config.Configuration {
			return cfg
		}),

		inject.Provide(db.NewPostgresConn),
		inject.Provide(db.NewRedisConn),
		inject.Provide(func() infra.Passport {
			return infra.NewPassport(12)
		}),

		inject.Provide(user.NewRepository),
		inject.Provide(todo.NewRepository),
		inject.Provide(auth.NewRepository),

		inject.Provide(auth.NewService),
		inject.Provide(func() auth.CreateAccessTokenHandler {
			return auth.CreateAccessToken
		}),
		inject.Provide(func() auth.CreateRefreshTokenHandler {
			return auth.CreateRefreshToken
		}),

		inject.Provide(common.NewController, inject.As(api.ControllerToken)),
		inject.Provide(user.NewController, inject.As(api.ControllerToken)),
		inject.Provide(todo.NewController, inject.As(api.ControllerToken)),
		inject.Provide(auth.NewController, inject.As(api.ControllerToken)),
	)

	return container, err
}
