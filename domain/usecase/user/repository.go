package user

import "github.com/gghcode/apas-todo-apiserver/domain/entity"

// Repository godoc
type Repository interface {
	CreateUser(User) (User, error)
	UserByID(userID int64) (User, error)
	UserByUserName(username string) (entity.User, error)
	UpdateUserByID(user User) (User, error)
	RemoveUserByID(userID int64) (User, error)
}
