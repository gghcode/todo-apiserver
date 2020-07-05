package bcrypt

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"

	"golang.org/x/crypto/bcrypt"
)

type passwordEncryptor struct {
	cost int
}

func NewPasswordEncryptor(cfg config.Configuration) user.PasswordEncryptor {
	return &passwordEncryptor{
		cost: cfg.BcryptCost,
	}
}

func (p *passwordEncryptor) HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
