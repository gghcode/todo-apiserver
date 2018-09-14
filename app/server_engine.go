package app

import "github.com/gin-gonic/gin"

type ServerEngine interface {
	Run(addr string) error
}

type GinEngine struct {
	gin *gin.Engine
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

