package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-gonic/gin"
)

// JwtAuthHandlerToken godoc
const JwtAuthHandlerToken = "JWT_AUTH_HANDLER_TOKEN"

// AddJwtAuthHandler godoc
func AddJwtAuthHandler(conf config.JwtConfig, authHandler *gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(JwtAuthHandlerToken, authHandler)
		ctx.Next()
	}
}

// JwtAuthRequired godoc
func JwtAuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyHandler := ctx.MustGet(JwtAuthHandlerToken).(*gin.HandlerFunc)
		(*verifyHandler)(ctx)

		ctx.Next()
	}
}
