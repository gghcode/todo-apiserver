package app

import (
	"path"

	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/resource/todo"
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

	registerRouter(server.core, "/api/v1/todo", todo.NewV1Router())

	return &server
}

func registerRouter(core *gin.Engine, basePath string, router api.Router) {
	for _, route := range router.Routes() {
		core.Handle(route.Method, path.Join(basePath, route.Path), route.Handler)
	}
}

// Run start listen.
func (server Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}
