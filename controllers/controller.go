package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	InitializeRoutes(router *gin.Engine)
}
