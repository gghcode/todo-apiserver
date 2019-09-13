package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api/user"
)

// UserRepository godoc
type UserRepository struct {
	mock.Mock
}

// CreateUser godoc
func (repo *UserRepository) CreateUser(usr user.User) (user.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(user.User), args.Error(1)
}

// AllUsers godoc
func (repo *UserRepository) AllUsers() ([]user.User, error) {
	args := repo.Called()
	return args.Get(0).([]user.User), args.Error(1)
}

// UserByID godoc
func (repo *UserRepository) UserByID(userID int64) (user.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}

// UserByUserName godoc
func (repo *UserRepository) UserByUserName(username string) (user.User, error) {
	args := repo.Called(username)
	return args.Get(0).(user.User), args.Error(1)
}

// UpdateUserByID godoc
func (repo *UserRepository) UpdateUserByID(usr user.User) (user.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(user.User), args.Error(1)
}

// RemoveUserByID godoc
func (repo *UserRepository) RemoveUserByID(userID int64) (user.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}