package security_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infra/security"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPassportPasswordVerification(t *testing.T) {
	testCases := []struct {
		description       string
		argBcryptCost     int
		argOriginPassword string
		argTargetPassword string
		expectedErr       error
		expected          bool
	}{
		{
			description:       "ShouldReturnErr",
			argBcryptCost:     bcrypt.MaxCost + 1,
			argOriginPassword: "1234",
			argTargetPassword: "12",
			expectedErr:       bcrypt.InvalidCostError(bcrypt.MaxCost + 1),
			expected:          false,
		},
		{
			description:       "ShouldBeValid",
			argBcryptCost:     bcrypt.MinCost,
			argOriginPassword: "12345678",
			argTargetPassword: "12345678",
			expected:          true,
		},
		{
			description:       "ShouldBeInvalid",
			argBcryptCost:     bcrypt.MinCost,
			argOriginPassword: "12345678910",
			argTargetPassword: "12345",
			expected:          false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			bcryptPassport := security.NewBcryptPassport(config.Configuration{
				BcryptCost: tc.argBcryptCost,
			})

			passwordHash, actualErr := bcryptPassport.HashPassword(tc.argOriginPassword)
			assert.Equal(t, tc.expectedErr, actualErr)

			actual := bcryptPassport.IsValidPassword(tc.argTargetPassword, passwordHash)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
