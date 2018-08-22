package controllers

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	GetHandlers() []Handler
}

type HandlerFunc = gin.HandlerFunc

type Handler struct {
	Method string
	Path   string
	Handle HandlerFunc
}
