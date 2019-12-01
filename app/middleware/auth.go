package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gin-gonic/gin"
)

// accessTokenHandlerToken godoc
const accessTokenHandlerToken = "ACCESS_TOKEN_HANDLER_TOKEN"

// AccessTokenHandlerFunc is function that handle access token
type AccessTokenHandlerFunc func(ctx *gin.Context) error

// AccessTokenHandlerFactory return AccessTokenHandlerFunc
type AccessTokenHandlerFactory interface {
	Create() AccessTokenHandlerFunc
}

// AddAccessTokenHandler godoc
func AddAccessTokenHandler(accessTokenHandlerFactory AccessTokenHandlerFactory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(accessTokenHandlerToken, accessTokenHandlerFactory.Create())
		ctx.Next()
	}
}

// RequiredAccessToken godoc
func RequiredAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessTokenHandler := ctx.MustGet(accessTokenHandlerToken).(AccessTokenHandlerFunc)
		if err := accessTokenHandler(ctx); err != nil {
			api.AbortErrorResponse(ctx, err)
			return
		}

		ctx.Next()
	}
}
