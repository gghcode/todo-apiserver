package auth

import "github.com/gghcode/apas-todo-apiserver/domain/entity"

// UserDataSource retrieves user data from detail source
type UserDataSource interface {
	UserByUserName(username string) (entity.User, error)
}
