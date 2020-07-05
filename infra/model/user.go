package model

import "github.com/gghcode/apas-todo-apiserver/domain/usecase/user"

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    int64  `gorm:"not null;"`
}

// FromUserEntity create user data model from user entity model
func FromUserEntity(usr user.User) User {
	return User{
		ID:           usr.ID,
		UserName:     usr.UserName,
		PasswordHash: usr.PasswordHash,
		CreatedAt:    usr.CreatedAt,
	}
}

// ToUserEntity create user entity model from user data model
func ToUserEntity(usr User) user.User {
	return user.User{
		ID:           usr.ID,
		UserName:     usr.UserName,
		PasswordHash: usr.PasswordHash,
		CreatedAt:    usr.CreatedAt,
	}
}
