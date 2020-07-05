package model

import (
	"github.com/gghcode/apas-todo-apiserver/domain/entity"
)

// User is user data model
type User struct {
	ID           int64  `gorm:"primary_key;"`
	UserName     string `gorm:"unique;not null;"`
	PasswordHash []byte `gorm:"not null;"`
	CreatedAt    int64  `gorm:"not null;"`
}

// FromUserEntity create user data model from user entity model
func FromUserEntity(usr entity.User) User {
	return User{
		ID:           usr.ID,
		UserName:     usr.UserName,
		PasswordHash: usr.PasswordHash,
		CreatedAt:    usr.CreatedAt,
	}
}

// ToUserEntity create user entity model from user data model
func ToUserEntity(usr User) entity.User {
	return entity.User{
		ID:           usr.ID,
		UserName:     usr.UserName,
		PasswordHash: usr.PasswordHash,
		CreatedAt:    usr.CreatedAt,
	}
}
