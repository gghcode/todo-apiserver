package auth

import (
	"time"

	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/infra"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
)

// TokenResponse godoc
type TokenResponse struct {
	Type         string `json:"type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

// Service godoc
type Service interface {
	IssueToken(req LoginRequest) (TokenResponse, error)
}

// NewService return new auth authService instance.
func NewService(
	conf config.Configuration,
	passport infra.Passport,
	userRepo user.Repository) Service {

	return &authService{
		secretKeyBytes:      []byte(conf.Jwt.SecretKey),
		accessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
		refreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		userRepo:            userRepo,
		passport:            passport,
	}
}

type authService struct {
	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	userRepo user.Repository
	passport infra.Passport
}

// IssueToken godoc
func (service *authService) IssueToken(req LoginRequest) (TokenResponse, error) {
	var userID int64
	if err := service.authenticate(req, &userID); err != nil {
		return TokenResponse{}, err
	}

	return TokenResponse{}, nil
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
