package auth

import (
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gin-gonic/gin"
)

type jwtAccessTokenHandler struct {
	jwtCfg config.JwtConfig
}

func (handler *jwtAccessTokenHandler) Create() middleware.AccessTokenHandlerFunc {
	return func(ctx *gin.Context) error {
		token := ctx.GetHeader("Authorization")
		tokenSecretKey := handler.jwtCfg.SecretKey

		claims, err := VerifyAccessToken(tokenSecretKey, token)
		if err != nil {
			return err
		}

		userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			return err
		}

		ctx.Set("user_id", userID)

		return nil
	}
}

// NewJwtAccessTokenHandlerFactory return new jwtAccessTokenHandler instance
func NewJwtAccessTokenHandlerFactory(cfg config.Configuration) middleware.AccessTokenHandlerFactory {
	return &jwtAccessTokenHandler{
		jwtCfg: cfg.Jwt,
	}
}

// verifyAccessToken godoc
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

// extractTokenClaims godoc
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
