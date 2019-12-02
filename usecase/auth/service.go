package auth

import (
	"github.com/gghcode/apas-todo-apiserver/app/api/user"
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/usecase/security"
)

// CreateAccessTokenHandler godoc
type CreateAccessTokenHandler func(userID int64) (string, error)

// CreateAccessTokenHandlerFactory is factory that return CreateAccessTokenHandler
type CreateAccessTokenHandlerFactory func(config.Configuration) CreateAccessTokenHandler

// CreateRefreshTokenHandler godoc
type CreateRefreshTokenHandler func(userID int64) (string, error)

// CreateRefreshTokenHandlerFactory is factory that return CreateRefreshTokenHandler
type CreateRefreshTokenHandlerFactory func(config.Configuration, TokenRepository) CreateRefreshTokenHandler

type authService struct {
	cfg                       config.JwtConfig
	tokenRepo                 TokenRepository
	userRepo                  user.Repository
	passport                  security.Passport
	createAccessTokenHandler  CreateAccessTokenHandler
	createRefreshTokenHandler CreateRefreshTokenHandler
}

// NewService return new auth authService instance.
func NewService(
	cfg config.Configuration,
	passport security.Passport,
	tokenRepo TokenRepository,
	userRepo user.Repository,
	accessTokenHandlerFactory CreateAccessTokenHandlerFactory,
	refreshTokenHandlerFactory CreateRefreshTokenHandlerFactory) AuthUsecaseInteractor {

	return &authService{
		cfg:                       cfg.Jwt,
		userRepo:                  userRepo,
		tokenRepo:                 tokenRepo,
		passport:                  passport,
		createAccessTokenHandler:  accessTokenHandlerFactory(cfg),
		createRefreshTokenHandler: refreshTokenHandlerFactory(cfg, tokenRepo),
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
