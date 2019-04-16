package todo

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	v1 "gitlab.com/gyuhwan/apas-todo-apiserver/app/resource/todo/v1"
)

// NewV1Router return new v1 router of user.
func NewV1Router() api.Router {
	controller := v1.Controller{}

	return v1.NewRouter(controller)
}
