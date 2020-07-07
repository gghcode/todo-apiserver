package auth

import (
	"github.com/gghcode/apas-todo-apiserver/config"
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
		userDataSource       UserDataSource
		authenticator        PasswordAuthenticator
		generateAccessToken  AccessTokenGeneratorFunc
		generateRefreshToken RefreshTokenGeneratorFunc
	}
)

// NewService return new auth authService instance.
func NewService(
	cfg config.Configuration,
	authenticator PasswordAuthenticator,
	tokenRepo TokenRepository,
	userDataSource UserDataSource,
	accessTokenGeneratorFunc AccessTokenGeneratorFunc,
	refreshTokenGeneratorFunc RefreshTokenGeneratorFunc) UseCase {

	return &authService{
		cfg:                  cfg,
		userDataSource:       userDataSource,
		tokenRepo:            tokenRepo,
		authenticator:        authenticator,
		generateAccessToken:  accessTokenGeneratorFunc,
		generateRefreshToken: refreshTokenGeneratorFunc,
	}
}

func (srv *authService) IssueToken(req LoginRequest) (TokenResponse, error) {
	var res TokenResponse

	var userID int64
	if err := srv.authenticate(req, &userID); err != nil {
		return res, err
	}

	accessToken, err := srv.generateAccessToken(userID)
	if err != nil {
		return res, err
	}

	refreshToken, err := srv.generateRefreshToken(userID)
	if err != nil {
		return res, err
	}

	res = TokenResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    srv.cfg.JwtAccessExpiresInSec,
	}

	return res, nil
}

func (srv *authService) RefreshToken(req AccessTokenByRefreshRequest) (TokenResponse, error) {
	var res TokenResponse

	userID, err := srv.tokenRepo.UserIDByRefreshToken(req.Token)
	if err != nil {
		return res, err
	}

	accessToken, err := srv.generateAccessToken(userID)
	if err != nil {
		return res, err
	}

	res = TokenResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		ExpiresIn:   srv.cfg.JwtAccessExpiresInSec,
	}

	return res, nil
}

func (srv *authService) authenticate(req LoginRequest, userID *int64) error {
	usr, err := srv.userDataSource.UserByUserName(req.Username)
	if err == user.ErrUserNotFound {
		return ErrInvalidCredential
	} else if err != nil {
		return err
	}

	if !srv.authenticator.IsValidPassword(req.Password, usr.PasswordHash) {
		return ErrInvalidCredential
	}

	*userID = usr.ID

	return nil
}
