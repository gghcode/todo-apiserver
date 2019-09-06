package middleware_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type JwtAuthMiddlewareUnit struct {
	suite.Suite

	fakeUserID      int64
	fakeAccessToken string
}

func TestJwtAuthMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(JwtAuthMiddlewareUnit))
}

func (suite *JwtAuthMiddlewareUnit) SetupTest() {
	suite.fakeUserID = 10

}

// func (suite *JwtAuthMiddlewareUnit) TestVerifyAccessToken() {
// 	testCases := []struct {
// 		description    string
// 		argSecret      string
// 		argAccessToken string
// 		expected       jwt.MapClaims
// 		expectedErr    error
// 	}{
// 		{
// 			description: "ShouldVerificationSuccess",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		suite.Run(tc.description, func() {
// 			actual, actualErr :=
// 		})
// 	}
// }
