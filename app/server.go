package app

import (
	"github.com/gin-gonic/gin"
)

// Server is api-server instance. it contains gin.Engine, middlewares, configuration.
type Server struct {
	core *gin.Engine
}

// Run start listen.
func (server Server) Run() error {

	return nil
}

// New return new server instance.
func New() *Server {
	return &Server{}
}
