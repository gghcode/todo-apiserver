package middleware_test

import (
	"net/http"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type (
	fakeAccessTokenHandlerFactory struct {
		errHolder *fakeErrHolder
	}
	fakeErrHolder struct {
		mock.Mock
	}
)

func (f *fakeAccessTokenHandlerFactory) Create() middleware.AccessTokenHandlerFunc {
	return func(token string) (middleware.TokenClaims, error) {
		return middleware.TokenClaims{}, f.errHolder.Error()
	}
}

func (h *fakeErrHolder) Error() error {
	args := h.Called()
	return args.Error(0)
}

func TestRequiredAccessToken(t *testing.T) {
	testCases := []struct {
		description    string
		stubErr        error
		expectedStatus int
	}{
		{
			description:    "ShouldReturnOK",
			stubErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			description:    "ShouldReturnUnauthorized",
			stubErr:        fake.ErrStub,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			errorHolder := &fakeErrHolder{}
			errorHolder.On("Error").Return(tc.stubErr)

			accessTokenHandler := &fakeAccessTokenHandlerFactory{
				errHolder: errorHolder,
			}

			ginRouter := gin.New()
			ginRouter.Use(middleware.AddAccessTokenHandler(accessTokenHandler))
			ginRouter.Use(middleware.RequiredAccessToken())
			ginRouter.GET("", func(ctx *gin.Context) {})

			actual := testutil.Response(
				t,
				ginRouter,
				"GET",
				"",
				nil,
			)

			assert.Equal(t, tc.expectedStatus, actual.StatusCode)
		})
	}
}

// func TestAuthUserID(t *testing.T) {
// 	var stubUserID int64 = 5

// 	expectedUserID := strconv.FormatInt(stubUserID, 10)

// 	ginRouter := gin.New()
// 	ginRouter.GET("", func(ctx *gin.Context) {
// 		middleware.SetAuthUserID(ctx, stubUserID)

// 		userID := middleware.AuthUserID(ctx)

// 		ctx.String(http.StatusOK, strconv.FormatInt(userID, 10))
// 	})

// 	actual := testutil.Response(
// 		t,
// 		ginRouter,
// 		"GET",
// 		"",
// 		nil,
// 	)

// 	assert.Equal(t, http.StatusOK, actual.StatusCode)

// 	actualUserID := testutil.StringFromIOReader(t, actual.Body)

// 	assert.Equal(t, expectedUserID, actualUserID)
// }
