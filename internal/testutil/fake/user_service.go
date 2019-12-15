package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/user"

	"github.com/stretchr/testify/mock"
)

type UserService struct {
	mock.Mock
}

// NewUserService return fake user service
func NewUserService() *UserService {
	return &UserService{}
}

func (srv *UserService) CreateUser(req user.CreateUserRequest) (user.UserResponse, error) {
	args := srv.Called(req)
	return args.Get(0).(user.UserResponse), args.Error(1)
}

func (srv *UserService) GetUserByUserID(userID int64) (user.UserResponse, error) {
	args := srv.Called(userID)
	return args.Get(0).(user.UserResponse), args.Error(1)
}

func (srv *UserService) GetUserByUserName(userName string) (user.UserResponse, error) {
	args := srv.Called(userName)
	return args.Get(0).(user.UserResponse), args.Error(1)
}
