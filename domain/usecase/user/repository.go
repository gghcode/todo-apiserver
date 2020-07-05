package user

import "github.com/gghcode/apas-todo-apiserver/domain/entity"

// Repository godoc
type Repository interface {
	CreateUser(entity.User) (entity.User, error)
	UserByID(userID int64) (entity.User, error)
	UserByUserName(username string) (entity.User, error)
	UpdateUserByID(user entity.User) (entity.User, error)
	RemoveUserByID(userID int64) (entity.User, error)
}
