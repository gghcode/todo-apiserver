package service

import (
	"github.com/defval/inject"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/todo"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
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

		inject.Provide(user.NewRepository),
		inject.Provide(todo.NewRepository),

		inject.Provide(common.NewController, inject.As(api.ControllerToken)),
		inject.Provide(user.NewController, inject.As(api.ControllerToken)),
		inject.Provide(todo.NewController, inject.As(api.ControllerToken)),
	)

	return container, err
}
