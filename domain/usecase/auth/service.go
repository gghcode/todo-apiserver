package auth

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/security"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
)

type (
	// AccessTokenGeneratorFunc generate access token
	AccessTokenGeneratorFunc func(userID int64) (string, error)

	// RefreshTokenGeneratorFunc generate refresh token
	RefreshTokenGeneratorFunc func(userID int64) (string, error)

	authService struct {
		cfg                  config.Configuration
		tokenRepo            TokenRepository
		userRepo             user.Repository
		passport             security.Passport
		generateAccessToken  AccessTokenGeneratorFunc
		generateRefreshToken RefreshTokenGeneratorFunc
	}
)

// NewService return new auth authService instance.
func NewService(
	cfg config.Configuration,
	passport security.Passport,
	tokenRepo TokenRepository,
	userRepo user.Repository,
	accessTokenGeneratorFunc AccessTokenGeneratorFunc,
	refreshTokenGeneratorFunc RefreshTokenGeneratorFunc) UseCase {

	return &authService{
		cfg:                  cfg,
		userRepo:             userRepo,
		tokenRepo:            tokenRepo,
		passport:             passport,
		generateAccessToken:  accessTokenGeneratorFunc,
		generateRefreshToken: refreshTokenGeneratorFunc,
	}
}

func (service *authService) IssueToken(req LoginRequest) (TokenResponse, error) {
	var res TokenResponse

	var userID int64
	if err := service.authenticate(req, &userID); err != nil {
		return res, err
	}

	accessToken, err := service.generateAccessToken(userID)
	if err != nil {
		return res, err
	}

	refreshToken, err := service.generateRefreshToken(userID)
	if err != nil {
		return res, err
	}

	res = TokenResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    service.cfg.JwtAccessExpiresInSec,
	}

	return res, nil
}

func (service *authService) RefreshToken(req AccessTokenByRefreshRequest) (TokenResponse, error) {
	var res TokenResponse

	userID, err := service.tokenRepo.UserIDByRefreshToken(req.Token)
	if err != nil {
		return res, err
	}

	accessToken, err := service.generateAccessToken(userID)
	if err != nil {
		return res, err
	}

	res = TokenResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		ExpiresIn:   service.cfg.JwtAccessExpiresInSec,
	}

	return res, nil
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
