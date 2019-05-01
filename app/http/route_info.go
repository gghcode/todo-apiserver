package http

import "github.com/gin-gonic/gin"

// RouteInfo include infomation about api route
type RouteInfo struct {
	Method  string
	Path    string
	Handler func(*gin.Context)
}
