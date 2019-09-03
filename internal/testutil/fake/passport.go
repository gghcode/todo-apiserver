package fake

import "github.com/stretchr/testify/mock"

// Passport godoc
type Passport struct {
	mock.Mock
}

// HashPassword godoc
func (passport *Passport) HashPassword(password string) ([]byte, error) {
	args := passport.Called(password)
	return args.Get(0).([]byte), args.Error(1)
}

// IsValidPassword godoc
func (passport *Passport) IsValidPassword(password string, hash []byte) bool {
	args := passport.Called(password, hash)
	return args.Bool(0)
}
