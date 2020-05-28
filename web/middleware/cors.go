package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware godoc
type CorsMiddleware gin.HandlerFunc

// NewCors godoc
func NewCors(cfg config.Configuration) CorsMiddleware {
	return CorsMiddleware(cors.New(cors.Config{
		AllowOrigins: cfg.CorsAllowOrigins,
		AllowMethods: cfg.CorsAllowMethods,
		AllowHeaders: []string{
			"Authorization",
			"Accept",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
			"Origin",
			"Content-Type",
			"X-Requested-With",
		},
		AllowCredentials: true,
	}))
}
