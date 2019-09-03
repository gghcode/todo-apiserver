package auth

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

// JWTParam godoc
type JWTParam struct {
	SecretKeyBytes      []byte
	AccessExpiresInSec  time.Duration
	RefreshExpiresInSec time.Duration
}

// EmptyTokenResponse godoc
var EmptyTokenResponse = TokenResponse{}

// CreateAccessTokenHandler godoc
type CreateAccessTokenHandler func(param JWTParam, userID int64) (string, error)

// CreateRefreshTokenHandler godoc
type CreateRefreshTokenHandler func(param JWTParam, userID int64) (string, error)

// Service godoc
type Service interface {
	IssueToken(req LoginRequest) (TokenResponse, error)
}

// NewService return new auth authService instance.
func NewService(
	conf config.Configuration,
	passport infra.Passport,
	userRepo user.Repository,
	createAccessTokenHandler CreateAccessTokenHandler,
	createRefreshTokenHandler CreateRefreshTokenHandler) Service {

	return &authService{
		cfg: conf,
		jwtParam: JWTParam{
			SecretKeyBytes:      []byte(conf.Jwt.SecretKey),
			AccessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
			RefreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		},
		userRepo:                  userRepo,
		passport:                  passport,
		createAccessTokenHandler:  createAccessTokenHandler,
		createRefreshTokenHandler: createRefreshTokenHandler,
	}
}

type authService struct {
	cfg      config.Configuration
	jwtParam JWTParam

	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	userRepo user.Repository
	passport infra.Passport

	createAccessTokenHandler  CreateAccessTokenHandler
	createRefreshTokenHandler CreateRefreshTokenHandler
}

// IssueToken godoc
func (service *authService) IssueToken(req LoginRequest) (TokenResponse, error) {
	var userID int64
	if err := service.authenticate(req, &userID); err != nil {
		return TokenResponse{}, err
	}

	accessToken, err := service.createAccessTokenHandler(service.jwtParam, userID)
	if err != nil {
		return EmptyTokenResponse, err
	}

	refreshToken, err := service.createRefreshTokenHandler(service.jwtParam, userID)
	if err != nil {
		return EmptyTokenResponse, err
	}

	return TokenResponse{
		Type:         "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    service.cfg.Jwt.AccessExpiresInSec,
	}, nil
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

// CreateAccessToken godoc
func CreateAccessToken(jwtParam JWTParam, userID int64) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   strconv.FormatInt(userID, 10),
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

// CreateRefreshToken godoc
func CreateRefreshToken(jwtParam JWTParam, userID int64) (string, error) {
	return "fasdf", nil
}

// ExtractTokenClaims godoc
func ExtractTokenClaims(jwtParam JWTParam, token string) (jwt.MapClaims, error) {
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
			return nil, err
		}

		return nil, err
	}

	return claims, nil
}
