package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetHandlerInfos() []HandlerInfo
}

type HandlerFunc = gin.HandlerFunc

type HandlerInfo struct {
	Method string
	Path   string
	Handle HandlerFunc
}
