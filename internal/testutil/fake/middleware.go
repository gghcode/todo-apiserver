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

// AddJwtAuthHandler godoc
func AddJwtAuthHandler(userIDFactory UserIDFactory) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var innerHandler gin.HandlerFunc = func(ctx *gin.Context) {
			ctx.Set("user_id", userIDFactory.UserID())
			ctx.Next()
		}

		ctx.Set(middleware.JwtAuthHandlerToken, &innerHandler)
		ctx.Next()
	}
}
