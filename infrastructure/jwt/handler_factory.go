package jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/auth"
)

type jwtParam struct {
	SecretKeyBytes      []byte
	AccessExpiresInSec  time.Duration
	RefreshExpiresInSec time.Duration
}

func NewJwtAccessTokenHandlerFactory(cfg config.Configuration) auth.CreateAccessTokenHandlerFactory {
	return CreateAccessTokenFactory
}

func NewJwtRefreshTokenHandlerfactory(cfg config.Configuration, tokenRepo auth.TokenRepository) auth.CreateRefreshTokenHandlerFactory {
	return CreateRefreshTokenFactory
}

// CreateAccessTokenFactory godoc
func CreateAccessTokenFactory(cfg config.Configuration) auth.CreateAccessTokenHandler {
	params := jwtParam{
		SecretKeyBytes:      []byte(cfg.Jwt.SecretKey),
		AccessExpiresInSec:  time.Duration(cfg.Jwt.AccessExpiresInSec),
		RefreshExpiresInSec: time.Duration(cfg.Jwt.RefreshExpiresInSec),
	}

	return func(userID int64) (string, error) {
		return createJwtToken(params, "access", userID)
	}
}

// CreateRefreshTokenFactory godoc
func CreateRefreshTokenFactory(cfg config.Configuration, tokenRepo auth.TokenRepository) auth.CreateRefreshTokenHandler {
	params := jwtParam{
		SecretKeyBytes:      []byte(cfg.Jwt.SecretKey),
		AccessExpiresInSec:  time.Duration(cfg.Jwt.AccessExpiresInSec),
		RefreshExpiresInSec: time.Duration(cfg.Jwt.RefreshExpiresInSec),
	}

	return func(userID int64) (string, error) {
		token, err := createJwtToken(params, "refresh", userID)
		if err != nil {
			return "", err
		}

		tokenRepo.SaveRefreshToken(
			userID,
			token,
			params.RefreshExpiresInSec,
		)

		return token, nil
	}
}

func createJwtToken(jwtParam jwtParam, tokenType string, sub int64) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(sub, 10),
		ExpiresAt: time.Now().Add(jwtParam.AccessExpiresInSec * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
	}

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenObj.SignedString(jwtParam.SecretKeyBytes)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
