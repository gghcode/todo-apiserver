package infra

import "golang.org/x/crypto/bcrypt"

// Passport is object that execute about password auth
type Passport interface {
	HashPassword(password string) ([]byte, error)
	IsValidPassword(password string, hash []byte) bool
}

type passport struct {
	cost int
}

func (passport *passport) HashPassword(password string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passport.cost)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (passport *passport) IsValidPassword(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

// NewPassport return new passport.
func NewPassport(cost int) Passport {
	return &passport{
		cost: cost,
	}
}
