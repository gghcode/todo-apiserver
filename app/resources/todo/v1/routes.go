package v1

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/http"
)

func Routes(c *controller) []http.RouteInfo {
	return []http.RouteInfo{
		http.RouteInfo{
			Method:  "GET",
			Path:    "/",
			Handler: c.getHandler,
		},
	}
}
