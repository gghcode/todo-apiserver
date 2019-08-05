package service

import (
	"github.com/defval/inject"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/common"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

// NewContainer godoc
func NewContainer(cfg config.Configuration) (*inject.Container, error) {
	container, err := inject.New(
		inject.Provide(cfg),

		inject.Provide(common.NewController(), inject.As(api.ControllerToken)),
	)

	return container, err
}
