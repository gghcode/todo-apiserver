package middleware

import (
	"strconv"

	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gghcode/apas-todo-apiserver/app/api/auth"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/app/val"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-gonic/gin"
)

// JwtAuthHandlerToken godoc
const JwtAuthHandlerToken = "JWT_AUTH_HANDLER_TOKEN"

// AddJwtAuthHandler godoc
func AddJwtAuthHandler(conf config.JwtConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var innerHandler gin.HandlerFunc = func(ctx *gin.Context) {
			token := ctx.GetHeader("Authorization")

			claims, err := auth.VerifyAccessToken(conf.SecretKey, token)
			if err != nil {
				api.AbortErrorResponse(ctx, err)
				return
			}

			userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
			if err != nil {
				api.AbortErrorResponse(ctx, user.ErrInvalidUserID)
				return
			}

			ctx.Set(val.UserID, userID)
			ctx.Next()
		}

		ctx.Set(JwtAuthHandlerToken, innerHandler)
		ctx.Next()
	}
}

// JwtAuthRequired godoc
func JwtAuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyHandler := ctx.MustGet(JwtAuthHandlerToken).(gin.HandlerFunc)
		verifyHandler(ctx)

		ctx.Next()
	}
}
