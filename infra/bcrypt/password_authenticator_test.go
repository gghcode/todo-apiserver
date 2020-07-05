package bcrypt_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/infra/bcrypt"
)

func TestBcryptPasswordAuthenticator(t * testing.T) {
	testCases := []struct {
		description string
	}{
		{
			description: "",
		}
	}

	for tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			authenticator := bcrypt.NewPasswordAuthenticator()
		})
	}
}
