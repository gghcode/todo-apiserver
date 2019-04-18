package app

import (
	"path"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/http"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/resources/todo"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

// Server is api-server instance. it contains gin.Engine, middlewares, configuration.
type Server struct {
	core *gin.Engine
	conf config.Configuration
}

// New return new server instance.
func New(conf config.Configuration) *Server {
	server := Server{
		core: gin.New(),
		conf: conf,
	}

	registerRouter(server.core, "/api/v1/todo", todo.NewV1Resource())

	return &server
}

func registerRouter(core *gin.Engine, basePath string, routes []http.RouteInfo) {
	for _, route := range routes {
		core.Handle(route.Method, path.Join(basePath, route.Path), route.Handler)
	}
}

// Run start listen.
func (server Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}
