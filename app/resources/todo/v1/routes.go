package v1

import (
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/http"
)

func Routes(c *Controller) []http.RouteInfo {
	return []http.RouteInfo{
		http.RouteInfo{
			Method:  "GET",
			Path:    "/",
			Handler: c.getAllHandler,
		},
		http.RouteInfo{
			Method: "GET",
			Path: "/:id",
			Handler: c.getByIdHandler,
		},
		http.RouteInfo{
			Method: "POST",
			Path: "/",
			Handler: c.createHandler,
		},
		http.RouteInfo{
			Method: "DELETE",
			Path: "/:id",
			Handler: c.removeByIdHandler,
		},
	}
}
