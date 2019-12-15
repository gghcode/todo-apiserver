package middleware

import (
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

// accessTokenHandlerToken godoc
const accessTokenHandlerToken = "ACCESS_TOKEN_HANDLER_TOKEN"
const tokenClaimsToken = "TOKEN_CLAIMS_TOKEN"
const userIDToken = "USER_ID_TOKEN"

type (
	// TokenClaims is infomation that contain in access token
	TokenClaims struct {
		UserID int64
	}

	// AccessTokenHandlerFunc is function that handle access token
	AccessTokenHandlerFunc func(token string) (TokenClaims, error)

	// AccessTokenHandlerFactory return AccessTokenHandlerFunc
	AccessTokenHandlerFactory interface {
		Create() AccessTokenHandlerFunc
	}
)

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
		token := ctx.GetHeader("Authorization")
		accessTokenHandler := ctx.MustGet(accessTokenHandlerToken).(AccessTokenHandlerFunc)

		claims, err := accessTokenHandler(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.MakeErrorResponse(err))
			return
		}

		ctx.Set(tokenClaimsToken, claims)
		ctx.Next()
	}
}

// AuthUserID return authenticated user id
func AuthUserID(ctx *gin.Context) int64 {
	return ctx.MustGet(tokenClaimsToken).(TokenClaims).UserID
}
