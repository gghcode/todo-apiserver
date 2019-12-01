package fake

import (
	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gin-gonic/gin"
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

func (handler *accessTokenHandlerFactory) Create() middleware.AccessTokenHandlerFunc {
	return func(ctx *gin.Context) error {
		ctx.Set("user_id", handler.userIDFactory.UserID())
		return nil
	}
}

func NewAccessTokenHandlerFactory(userIDFactory UserIDFactory) middleware.AccessTokenHandlerFactory {
	return &accessTokenHandlerFactory{
		userIDFactory: userIDFactory,
	}
}
