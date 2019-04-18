package todo

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/http"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/resources/todo/v1"
)

// NewV1Resource returns v1 routes of user.
func NewV1Resource() []http.RouteInfo  {
	c := v1.NewController()

	return v1.Routes(c)
}
