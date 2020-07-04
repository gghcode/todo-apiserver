package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/stretchr/testify/mock"
)

// AuthService godoc
type AuthService struct {
	mock.Mock
}

// NewAuthService return fake auth service
func NewAuthService() *AuthService {
	return &AuthService{}
}

// IssueToken godoc
func (service *AuthService) IssueToken(req auth.LoginRequest) (auth.TokenResponse, error) {
	args := service.Called(req)
	return args.Get(0).(auth.TokenResponse), args.Error(1)
}

// RefreshToken godoc
func (service *AuthService) RefreshToken(req auth.AccessTokenByRefreshRequest) (auth.TokenResponse, error) {
	args := service.Called(req)
	return args.Get(0).(auth.TokenResponse), args.Error(1)
}
