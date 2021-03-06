package fake

import (
	"github.com/gghcode/apas-todo-apiserver/web/middleware"
	"github.com/stretchr/testify/mock"
)

// UserIDFactory godoc
type UserIDFactory interface {
	UserID() int64
}

// MockUserID godoc
type MockUserID struct {
	mock.Mock
}

// UserID godoc
func (m *MockUserID) UserID() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

type accessTokenHandlerFactory struct {
	userIDFactory UserIDFactory
}

// NewAccessTokenHandlerFactory return new fake access token handler factory
func NewAccessTokenHandlerFactory(userIDFactory UserIDFactory) middleware.AccessTokenHandlerFactory {
	return &accessTokenHandlerFactory{
		userIDFactory: userIDFactory,
	}
}

func (handler *accessTokenHandlerFactory) Create() middleware.AccessTokenHandlerFunc {
	return func(token string) (middleware.TokenClaims, error) {
		return middleware.TokenClaims{
			UserID: handler.userIDFactory.UserID(),
		}, nil
	}
}
