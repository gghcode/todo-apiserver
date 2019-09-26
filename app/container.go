package app

import (
	"github.com/defval/inject"
	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/auth"
	"github.com/gghcode/apas-todo-apiserver/app/api/common"
	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/app/infra"
	"github.com/gghcode/apas-todo-apiserver/app/loader"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/spf13/afero"
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
		inject.Provide(func() afero.Fs {
			return afero.NewOsFs()
		}),

		inject.Provide(loader.NewVersionLoader),

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
