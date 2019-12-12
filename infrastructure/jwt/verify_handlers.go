package jwt

import (
	"strconv"
	"strings"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type jwtAccessTokenHandler struct {
	jwtCfg config.JwtConfig
}

// NewJwtAccessTokenVerifyHandlerFactory return new jwtAccessTokenHandler instance
func NewJwtAccessTokenVerifyHandlerFactory(cfg config.Configuration) middleware.AccessTokenHandlerFactory {
	return &jwtAccessTokenHandler{
		jwtCfg: cfg.Jwt,
	}
}

func (handler *jwtAccessTokenHandler) Create() middleware.AccessTokenHandlerFunc {
	secretKeyBytes := []byte(handler.jwtCfg.SecretKey)

	return func(ctx *gin.Context) error {
		token := ctx.GetHeader("Authorization")

		claims, err := verifyAccessToken(secretKeyBytes, token)
		if err != nil {
			return err
		}

		userID, err := strconv.ParseInt(claims["sub"].(string), 10, 64)
		if err != nil {
			return err
		}

		middleware.SetAuthUserID(ctx, userID)

		return nil
	}
}

func verifyAccessToken(secretKeyBytes []byte, accessToken string) (jwt.MapClaims, error) {
	if accessToken == "" {
		return nil, auth.ErrNotContainTokenInHeader
	}

	tokenInfo := strings.Split(accessToken, " ")
	if len(tokenInfo) != 2 {
		return nil, auth.ErrInvalidToken
	}

	tokenType := tokenInfo[0]
	tokenString := tokenInfo[1]

	if tokenType != "Bearer" {
		return nil, auth.ErrInvalidTokenType
	}

	return extractTokenClaims(secretKeyBytes, tokenString)
}

// extractTokenClaims godoc
func extractTokenClaims(secretKeyBytes []byte, token string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	_, err := jwt.ParseWithClaims(
		token,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return secretKeyBytes, nil
		},
	)

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)

		if ok && validationErr.Errors == jwt.ValidationErrorExpired {
			return nil, auth.ErrTokenExpired
		}

		return nil, auth.ErrInvalidToken
	}

	return claims, nil
}