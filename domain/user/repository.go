package user

import "github.com/gghcode/apas-todo-apiserver/domain/model"

// Repository godoc
type Repository interface {
	CreateUser(model.User) (model.User, error)
	AllUsers() ([]model.User, error)
	UserByID(userID int64) (model.User, error)
	UserByUserName(username string) (model.User, error)
	UpdateUserByID(user model.User) (model.User, error)
	RemoveUserByID(userID int64) (model.User, error)
}
