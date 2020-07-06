package bcrypt

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordEncryptor(t *testing.T) {
	testCases := []struct {
		description   string
		argBcryptCost int
		argPassword   string
		expectedErr   error
	}{
		{
			description:   "ShouldReturnInvalidCostErr",
			argBcryptCost: bcrypt.MaxCost + 1,
			argPassword:   "testtest",
			expectedErr:   bcrypt.InvalidCostError(bcrypt.MaxCost + 1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			e := NewPasswordEncryptor(config.Configuration{
				BcryptCost: tc.argBcryptCost,
			})

			_, actualErr := e.HashPassword(tc.argPassword)

			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
