package app

import (
	"apas-todo-apiserver/controllers"
	"github.com/gin-gonic/gin"
)

type ServerEngine interface {
	Handle(handlerInfo controllers.HandlerInfo)
	Run(addr string) error
}

type GinEngine struct {
	gin *gin.Engine
}

func (engine *GinEngine) Handle(info controllers.HandlerInfo) {
	engine.gin.Handle(info.Method, info.Path, info.Handle)
}

func (engine *GinEngine) Run(addr string) error {
	return engine.gin.Run(addr)
}

func NewGinEngine() *GinEngine {
	result := GinEngine{
		gin: gin.New(),
	}

	return &result
}

