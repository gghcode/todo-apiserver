package infra

import "golang.org/x/crypto/bcrypt"

// Passport is object that execute about password auth
type Passport interface {
	HashPassword(password string) ([]byte, error)
	IsValidPassword(password string, hash []byte) bool
}

type bcryptPassport struct {
	cost int
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

// NewBcryptPassport return new passport.
func NewBcryptPassport(cost int) Passport {
	return &bcryptPassport{
		cost: cost,
	}
}
