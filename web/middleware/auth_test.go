package middleware_test

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/api"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	fakeAccessTokenHandlerFactory struct {
		stubHandler *fakeHandler
	}
	fakeHandler struct {
		mock.Mock
	}
)

func (f *fakeAccessTokenHandlerFactory) Create() middleware.AccessTokenHandlerFunc {
	return func(token string) (middleware.TokenClaims, error) {
		return f.stubHandler.Pass()
	}
}

func (h *fakeHandler) Pass() (middleware.TokenClaims, error) {
	args := h.Called()
	return args.Get(0).(middleware.TokenClaims), args.Error(1)
}

func TestRequiredAccessToken(t *testing.T) {
	testCases := []struct {
		description     string
		stubTokenClaims middleware.TokenClaims
		stubErr         error
		expectedStatus  int
		expectedJSON    string
	}{
		{
			description: "ShouldReturnOK",
			stubTokenClaims: middleware.TokenClaims{
				UserID: 5,
			},
			stubErr:        nil,
			expectedStatus: http.StatusOK,
			expectedJSON:   "5",
		},
		{
			description:    "ShouldReturnUnauthorized",
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusUnauthorized,
			expectedJSON: testutil.JSONStringFromInterface(t,
				api.MakeErrorResponseDTO(fake.ErrStub),
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			errorHolder := &fakeHandler{}
			errorHolder.On("Pass").Return(tc.stubTokenClaims, tc.stubErr)

			accessTokenHandler := &fakeAccessTokenHandlerFactory{
				stubHandler: errorHolder,
			}

			ginRouter := gin.New()
			ginRouter.Use(middleware.AddAccessTokenHandler(accessTokenHandler))
			ginRouter.Use(middleware.RequiredAccessToken())
			ginRouter.GET("/", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, strconv.FormatInt(
					middleware.AuthUserID(ctx),
					10,
				))
			})

			actual := testutil.Response(
				t,
				ginRouter,
				"GET",
				"/",
				nil,
			)

			assert.Equal(t, tc.expectedStatus, actual.StatusCode)

			actualJSON := testutil.StringFromIOReader(t, actual.Body)

			assert.JSONEq(t, tc.expectedJSON, actualJSON)
		})
	}
}
