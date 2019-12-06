package app

import (
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
)

type appController struct {
}

// NewController return new app controller
func NewController() api.GinController {
	return &appController{}
}

func (c *appController) RegisterRoutes(router gin.IRouter) {

}
