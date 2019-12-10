package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// accessTokenHandlerToken godoc
const accessTokenHandlerToken = "ACCESS_TOKEN_HANDLER_TOKEN"
const userIDToken = "USER_ID_TOKEN"

// AccessTokenHandlerFunc is function that handle access token
type AccessTokenHandlerFunc func(ctx *gin.Context) error

// AccessTokenHandlerFactory return AccessTokenHandlerFunc
type AccessTokenHandlerFactory interface {
	Create() AccessTokenHandlerFunc
}

// AddAccessTokenHandler godoc
func AddAccessTokenHandler(accessTokenHandlerFactory AccessTokenHandlerFactory) gin.HandlerFunc {
	accessTokenHandler := accessTokenHandlerFactory.Create()

	return func(ctx *gin.Context) {
		ctx.Set(accessTokenHandlerToken, accessTokenHandler)
		ctx.Next()
	}
}

// RequiredAccessToken godoc
func RequiredAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessTokenHandler := ctx.MustGet(accessTokenHandlerToken).(AccessTokenHandlerFunc)
		if err := accessTokenHandler(ctx); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.MakeErrorResponse(err))
			return
		}

		ctx.Next()
	}
}

// SetAuthUserID set authenticated user id
func SetAuthUserID(ctx *gin.Context, userID int64) {
	ctx.Set(userIDToken, userID)
}

// AuthUserID return authenticated user id
func AuthUserID(ctx *gin.Context) int64 {
	return ctx.GetInt64(userIDToken)
}
