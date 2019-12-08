package web

import (
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
)

// NewMux return mux
func NewMux(controllers []api.GinController) *gin.Engine {
	router := gin.New()

	for _, c := range controllers {
		c.RegisterRoutes(router)
	}

	return router
}
