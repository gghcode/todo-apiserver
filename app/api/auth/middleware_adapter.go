package auth

import (
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/app/api"
	"github.com/gin-gonic/gin"
)

// JwtAuthHandler godoc
var JwtAuthHandler gin.HandlerFunc = func(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	tokenSecretKey := ctx.GetString("tok_secret")

	claims, err := VerifyAccessToken(tokenSecretKey, token)
	if err != nil {
		api.AbortErrorResponse(ctx, err)
		return
	}

	userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
	if err != nil {
		api.AbortErrorResponse(ctx, ErrInvalidToken)
		return
	}

	ctx.Set("user_id", userID)
	ctx.Next()
}

// VerifyAccessToken godoc
func VerifyAccessToken(secret string, accessToken string) (jwt.MapClaims, error) {
	if accessToken == "" {
		return nil, ErrNotContainToken
	}

	tokenInfo := strings.Split(accessToken, " ")
	if len(tokenInfo) != 2 {
		return nil, ErrInvalidToken
	}

	tokenType := tokenInfo[0]
	tokenString := tokenInfo[1]

	if tokenType != "Bearer" {
		return nil, ErrInvalidTokenType
	}

	param := JwtParam{
		SecretKeyBytes: []byte(secret),
	}

	return ExtractTokenClaims(param, tokenString)
}

// ExtractTokenClaims godoc
func ExtractTokenClaims(jwtParam JwtParam, token string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtParam.SecretKeyBytes, nil
		},
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)

		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return nil, ErrTokenExpired
		}

		return nil, ErrInvalidToken
	}

	return claims, nil
}
