package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/stretchr/testify/mock"
)

// User2Repository godoc
type User2Repository struct {
	mock.Mock
}

// NewUserRepository return fake user repository
func NewUserRepository() *User2Repository {
	return &User2Repository{}
}

// CreateUser godoc
func (repo *User2Repository) CreateUser(usr user.User) (user.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(user.User), args.Error(1)
}

// AllUsers godoc
func (repo *User2Repository) AllUsers() ([]user.User, error) {
	args := repo.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

// UserByID godoc
func (repo *User2Repository) UserByID(userID int64) (user.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}

// UserByUserName godoc
func (repo *User2Repository) UserByUserName(username string) (user.User, error) {
	args := repo.Called(username)
	return args.Get(0).(user.User), args.Error(1)
}

// UpdateUserByID godoc
func (repo *User2Repository) UpdateUserByID(usr user.User) (user.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(user.User), args.Error(1)
}

// RemoveUserByID godoc
func (repo *User2Repository) RemoveUserByID(userID int64) (user.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}
