package fake

import (
	"github.com/gghcode/apas-todo-apiserver/domain/entity"
	"github.com/stretchr/testify/mock"
)

// UserRepository godoc
type UserRepository struct {
	mock.Mock
}

// NewUserRepository return fake user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CreateUser godoc
func (repo *UserRepository) CreateUser(usr entity.User) (entity.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(entity.User), args.Error(1)
}

// AllUsers godoc
func (repo *UserRepository) AllUsers() ([]entity.User, error) {
	args := repo.Called()
	return args.Get(0).([]entity.User), args.Error(1)
}

// UserByID godoc
func (repo *UserRepository) UserByID(userID int64) (entity.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(entity.User), args.Error(1)
}

// UserByUserName godoc
func (repo *UserRepository) UserByUserName(username string) (entity.User, error) {
	args := repo.Called(username)
	return args.Get(0).(entity.User), args.Error(1)
}

// UpdateUserByID godoc
func (repo *UserRepository) UpdateUserByID(usr entity.User) (entity.User, error) {
	args := repo.Called(usr)
	return args.Get(0).(entity.User), args.Error(1)
}

// RemoveUserByID godoc
func (repo *UserRepository) RemoveUserByID(userID int64) (entity.User, error) {
	args := repo.Called(userID)
	return args.Get(0).(entity.User), args.Error(1)
}
