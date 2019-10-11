package infra_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/infra"
	"github.com/stretchr/testify/suite"
)

type PassportUnit struct {
	suite.Suite

	passport infra.Passport
}

func TestPassportUnit(t *testing.T) {
	suite.Run(t, new(PassportUnit))
}

func (suite *PassportUnit) SetupTest() {
	suite.passport = infra.NewPassport(1)
}

func (suite *PassportUnit) TestPasswordVerification() {
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
		suite.Run(tc.description, func() {
			passwordHash, err := suite.passport.HashPassword(tc.password)
			suite.NoError(err)

			actual := suite.passport.IsValidPassword(tc.verifyPassword, passwordHash)
			suite.Equal(tc.expected, actual)
		})
	}
}
