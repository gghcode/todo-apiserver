package security

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/security"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPassport struct {
	cost int
}

// NewBcryptPassport return new passport.
func NewBcryptPassport(cfg config.Configuration) security.Passport {
	return &bcryptPassport{
		cost: cfg.BcryptCost,
	}
}

func (passport *bcryptPassport) HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passport.cost)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (passport *bcryptPassport) IsValidPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}
