package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	AddRoutes(router *gin.Engine)
}
