package jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
)

type (
	jwtOptions struct {
		SecretKeyBytes []byte
		ExpiresInSec   time.Duration
	}

	jwtClaims struct {
		jwt.StandardClaims

		TokenType string `json:"type"`
	}
)

// NewJwtAccessTokenGeneratorFunc return new jwt access token generator function
func NewJwtAccessTokenGeneratorFunc(cfg config.Configuration) auth.AccessTokenGeneratorFunc {
	opt := jwtOptions{
		SecretKeyBytes: []byte(cfg.Jwt.SecretKey),
		ExpiresInSec:   time.Duration(cfg.Jwt.AccessExpiresInSec),
	}

	return func(userID int64) (string, error) {
		claims := jwtClaims{
			generateStandardClaims(opt, userID),
			"access", // TokenType
		}

		return createToken(opt, claims)
	}
}

// NewJwtRefreshTokenGeneratorFunc return new jwt refresh token generator function
func NewJwtRefreshTokenGeneratorFunc(cfg config.Configuration) auth.RefreshTokenGeneratorFunc {
	opt := jwtOptions{
		SecretKeyBytes: []byte(cfg.Jwt.SecretKey),
		ExpiresInSec:   time.Duration(cfg.Jwt.RefreshExpiresInSec),
	}

	return func(userID int64) (string, error) {
		claims := jwtClaims{
			generateStandardClaims(opt, userID),
			"refresh", // TokenType
		}

		return createToken(opt, claims)
	}
}

func generateStandardClaims(opt jwtOptions, userID int64) jwt.StandardClaims {
	return jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: time.Now().Add(opt.ExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
}

func createToken(opt jwtOptions, claims jwtClaims) (string, error) {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tok.SignedString(opt.SecretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
