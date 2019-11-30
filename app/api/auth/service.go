package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/app/infra"
	"github.com/gghcode/apas-todo-apiserver/config"
)

// JwtParam godoc
type JwtParam struct {
	SecretKeyBytes      []byte
	AccessExpiresInSec  time.Duration
	RefreshExpiresInSec time.Duration
}

// CreateAccessTokenHandler godoc
type CreateAccessTokenHandler func(userID int64) (string, error)

// CreateAccessTokenHandlerFactory is factory that return CreateAccessTokenHandler
type CreateAccessTokenHandlerFactory func(JwtParam) CreateAccessTokenHandler

// CreateRefreshTokenHandler godoc
type CreateRefreshTokenHandler func(userID int64) (string, error)

// CreateRefreshTokenHandlerFactory is factory that return CreateRefreshTokenHandler
type CreateRefreshTokenHandlerFactory func(JwtParam, Repository) CreateRefreshTokenHandler

type authService struct {
	cfg                       config.JwtConfig
	tokenRepo                 Repository
	userRepo                  user.Repository
	passport                  infra.Passport
	createAccessTokenHandler  CreateAccessTokenHandler
	createRefreshTokenHandler CreateRefreshTokenHandler
}

// NewService return new auth authService instance.
func NewService(
	cfg config.JwtConfig,
	passport infra.Passport,
	tokenRepo Repository,
	userRepo user.Repository,
	accessTokenHandlerFactory CreateAccessTokenHandlerFactory,
	refreshTokenHandlerFactory CreateRefreshTokenHandlerFactory) Service {

	params := JwtParam{
		SecretKeyBytes:      []byte(cfg.SecretKey),
		AccessExpiresInSec:  time.Duration(cfg.AccessExpiresInSec),
		RefreshExpiresInSec: time.Duration(cfg.RefreshExpiresInSec),
	}

	return &authService{
		cfg:                       cfg,
		userRepo:                  userRepo,
		tokenRepo:                 tokenRepo,
		passport:                  passport,
		createAccessTokenHandler:  accessTokenHandlerFactory(params),
		createRefreshTokenHandler: refreshTokenHandlerFactory(params, tokenRepo),
	}
}

// IssueToken godoc
func (service *authService) IssueToken(req LoginRequest, res *TokenResponse) error {
	var userID int64
	if err := service.authenticate(req, &userID); err != nil {
		return err
	}

	accessToken, err := service.createAccessTokenHandler(userID)
	if err != nil {
		return err
	}

	refreshToken, err := service.createRefreshTokenHandler(userID)
	if err != nil {
		return err
	}

	*res = TokenResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    service.cfg.AccessExpiresInSec,
	}

	return nil
}

func (service *authService) RefreshToken(req AccessTokenByRefreshRequest, res *TokenResponse) error {
	userID, err := service.tokenRepo.UserIDByRefreshToken(req.Token)
	if err != nil {
		return err
	}

	accessToken, err := service.createAccessTokenHandler(userID)
	if err != nil {
		return err
	}

	*res = TokenResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		ExpiresIn:   service.cfg.AccessExpiresInSec,
	}

	return nil
}

func (service *authService) authenticate(req LoginRequest, userID *int64) error {
	usr, err := service.userRepo.UserByUserName(req.Username)
	if err == user.ErrUserNotFound {
		return ErrInvalidCredential
	} else if err != nil {
		return err
	}

	if !service.passport.IsValidPassword(req.Password, usr.PasswordHash) {
		return ErrInvalidCredential
	}

	*userID = usr.ID

	return nil
}

// CreateAccessTokenFactory godoc
func CreateAccessTokenFactory(jwtParam JwtParam) CreateAccessTokenHandler {
	return func(userID int64) (string, error) {
		return createJwtToken(jwtParam, "access", userID)
	}
}

// CreateRefreshTokenFactory godoc
func CreateRefreshTokenFactory(jwtParam JwtParam, tokenRepo Repository) CreateRefreshTokenHandler {
	return func(userID int64) (string, error) {
		token, err := createJwtToken(jwtParam, "refresh", userID)
		if err != nil {
			return "", err
		}

		tokenRepo.SaveRefreshToken(
			userID,
			token,
			jwtParam.RefreshExpiresInSec,
		)

		return token, nil
	}
}

func createJwtToken(jwtParam JwtParam, tokenType string, sub int64) (string, error) {
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
