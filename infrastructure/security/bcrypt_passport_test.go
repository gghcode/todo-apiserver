package security_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/security"
	"github.com/stretchr/testify/assert"
)

func TestBcryptPassportPasswordVerification(t *testing.T) {
	testCases := []struct {
		description    string
		password       string
		verifyPassword string
		expected       bool
	}{
		{
			description:    "ShouldBeValid",
			password:       "12345678",
			verifyPassword: "12345678",
			expected:       true,
		},
		{
			description:    "ShouldBeInvalid",
			password:       "12345678910",
			verifyPassword: "12345",
			expected:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			bcryptPassport := security.NewBcryptPassport(config.Configuration{
				BcryptCost: 1,
			})

			passwordHash, err := bcryptPassport.HashPassword(tc.password)
			assert.NoError(t, err)

			actual := bcryptPassport.IsValidPassword(tc.verifyPassword, passwordHash)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
