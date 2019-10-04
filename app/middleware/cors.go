package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCors godoc
func NewCors(cfg config.CorsConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}
