package v1

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

// Router provide Routes method.
type Router struct {
	controller Controller
}

// NewRouter return new router instance.
func NewRouter(controller Controller) *Router {
	return &Router{controller: controller}
}

// Routes returns slice of RouteInfo.
func (router Router) Routes() []api.RouteInfo {
	return []api.RouteInfo{
		api.RouteInfo{
			Method:  "GET",
			Path:    "/",
			Handler: router.controller.getHandler,
		},
	}
}
