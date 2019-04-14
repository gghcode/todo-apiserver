package app

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

// Server is api-server instance. it contains gin.Engine, middlewares, configuration.
type Server struct {
	core          *gin.Engine
	configuration *config.Configuration
}

// Run start listen.
func (server Server) Run() error {
	addr := server.configuration.Addr
	return server.core.Run(addr)
}

// New return new server instance.
func New() *Server {
	return &Server{}
}
