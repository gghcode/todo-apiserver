package bcrypt

import (
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"golang.org/x/crypto/bcrypt"
)

type passwordAuthenticator struct{}

func NewPasswordAuthenticator() auth.PasswordAuthenticator {
	return &passwordAuthenticator{}
}

func (p *passwordAuthenticator) IsValidPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
