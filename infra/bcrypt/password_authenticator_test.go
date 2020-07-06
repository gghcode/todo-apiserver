package bcrypt

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordAuthenticator(t *testing.T) {
	stubPasswd := "testtest"
	stubPasswdHash, _ := bcrypt.GenerateFromPassword([]byte(stubPasswd), bcrypt.MinCost)

	testCases := []struct {
		description      string
		argPassword      string
		stubPasswordHash []byte
		expected         bool
	}{
		{
			description:      "ShouldBeValid",
			argPassword:      stubPasswd,
			stubPasswordHash: stubPasswdHash,
			expected:         true,
		},
		{
			description:      "ShouldBeInvalid",
			argPassword:      "INVALID_PASSWORD",
			stubPasswordHash: stubPasswdHash,
			expected:         false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			authenticator := NewPasswordAuthenticator()
			authenticator.IsValidPassword("fasdf", nil)
		})
	}
}
