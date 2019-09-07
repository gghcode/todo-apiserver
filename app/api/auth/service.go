package auth

import (
	"strconv"
	"strings"
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

// JwtParam godoc
type JwtParam struct {
	SecretKeyBytes      []byte
	AccessExpiresInSec  time.Duration
	RefreshExpiresInSec time.Duration
}

// EmptyTokenResponse godoc
var EmptyTokenResponse = TokenResponse{}

// CreateAccessTokenHandler godoc
type CreateAccessTokenHandler func(param JwtParam, userID int64) (string, error)

// CreateRefreshTokenHandler godoc
type CreateRefreshTokenHandler func(tokenRepo Repository, param JwtParam, userID int64) (string, error)

// Service godoc
type Service interface {
	IssueToken(req LoginRequest) (TokenResponse, error)
	RefreshToken(req AccessTokenByRefreshRequest) (TokenResponse, error)
}

// NewService return new auth authService instance.
func NewService(
	conf config.Configuration,
	passport infra.Passport,
	tokenRepo Repository,
	userRepo user.Repository,
	createAccessTokenHandler CreateAccessTokenHandler,
	createRefreshTokenHandler CreateRefreshTokenHandler) Service {

	return &authService{
		cfg: conf,
		jwtParam: JwtParam{
			SecretKeyBytes:      []byte(conf.Jwt.SecretKey),
			AccessExpiresInSec:  time.Duration(conf.Jwt.AccessExpiresInSec),
			RefreshExpiresInSec: time.Duration(conf.Jwt.RefreshExpiresInSec),
		},
		userRepo:                  userRepo,
		tokenRepo:                 tokenRepo,
		passport:                  passport,
		createAccessTokenHandler:  createAccessTokenHandler,
		createRefreshTokenHandler: createRefreshTokenHandler,
	}
}

type authService struct {
	cfg      config.Configuration
	jwtParam JwtParam

	secretKeyBytes      []byte
	accessExpiresInSec  time.Duration
	refreshExpiresInSec time.Duration

	tokenRepo Repository
	userRepo  user.Repository
	passport  infra.Passport

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

	refreshToken, err := service.createRefreshTokenHandler(service.tokenRepo, service.jwtParam, userID)
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

func (service *authService) RefreshToken(req AccessTokenByRefreshRequest) (TokenResponse, error) {
	userID, err := service.tokenRepo.UserIDByRefreshToken(req.Token)
	if err != nil {
		return EmptyTokenResponse, err
	}

	accessToken, err := service.createAccessTokenHandler(service.jwtParam, userID)
	if err != nil {
		return EmptyTokenResponse, err
	}

	return TokenResponse{
		Type:        "Bearer",
		AccessToken: accessToken,
		ExpiresIn:   service.cfg.Jwt.AccessExpiresInSec,
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
func CreateAccessToken(jwtParam JwtParam, userID int64) (string, error) {
	return createJwtToken(jwtParam, "access", userID)
}

// CreateRefreshToken godoc
func CreateRefreshToken(tokenRepo Repository, jwtParam JwtParam, userID int64) (string, error) {
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

// VerifyAccessToken godoc
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

// ExtractTokenClaims godoc
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
