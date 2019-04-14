package app

import (
	"github.com/gin-gonic/gin"
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

	return &server
}

// Run start listen.
func (server Server) Run() error {
	addr := server.conf.Addr
	return server.core.Run(addr)
}
