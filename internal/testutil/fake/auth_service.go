package fake

import (
	"github.com/gghcode/apas-todo-apiserver/app/api/auth"
	"github.com/stretchr/testify/mock"
)

// AuthService godoc
type AuthService struct {
	mock.Mock
}

// IssueToken godoc
func (service *AuthService) IssueToken(req auth.LoginRequest, res *auth.TokenResponse) error {
	args := service.Called(req, res)
	return args.Error(0)
}

// RefreshToken godoc
func (service *AuthService) RefreshToken(req auth.AccessTokenByRefreshRequest, res *auth.TokenResponse) error {
	args := service.Called(req, res)
	return args.Error(0)
}
