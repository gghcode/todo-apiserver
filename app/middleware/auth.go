package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-gonic/gin"
)

// accessTokenHandlerToken godoc
const accessTokenHandlerToken = "JWT_AUTH_HANDLER_TOKEN"

// AddAccessTokenHandler godoc
func AddAccessTokenHandler(conf config.JwtConfig, accessTokenHandler *gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("tok_secret", conf.SecretKey)
		ctx.Set(accessTokenHandlerToken, accessTokenHandler)
		ctx.Next()
	}
}

// RequiredAccessToken godoc
func RequiredAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessTokenHandler := ctx.MustGet(accessTokenHandlerToken).(*gin.HandlerFunc)
		(*accessTokenHandler)(ctx)

		ctx.Next()
	}
}
