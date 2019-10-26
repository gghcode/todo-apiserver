package todo_test

import (
	"io"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ValidatorUnit struct {
	suite.Suite

	router *gin.Engine
}

func TestTodoValidatorUnit(t *testing.T) {
	suite.Run(t, new(ValidatorUnit))
}

func (suite *ValidatorUnit) SetupTest() {
	suite.router = gin.New()

	gin.SetMode(gin.TestMode)
}

func (suite *ValidatorUnit) TestUpdateTodoBind() {
	var testFunc gin.HandlerFunc
	suite.router.POST("/test", func(ctx *gin.Context) {
		testFunc(ctx)
	})

	fakeTitle := "fake title"

	testCases := []struct {
		description       string
		reqPayload        io.Reader
		expectedErr       error
		expectedBindModel todo.UpdateTodoRequest
	}{
		{
			description: "ShouldBindModelWithTitle",
			reqPayload: testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
				"title": fakeTitle,
			}),
			expectedErr: nil,
			expectedBindModel: todo.UpdateTodoRequest{
				Title: &fakeTitle,
			},
		},
		{
			description: "ShouldReturnErrWhenNotContentsMin",
			reqPayload: testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
				"contents": 5,
			}),
			expectedErr:       nil,
			expectedBindModel: todo.UpdateTodoRequest{},
		},
	}

	for _, tc := range testCases {
		testFunc = func(ctx *gin.Context) {
			bodyValidator := todo.NewUpdateTodoValidator()
			actualErr := bodyValidator.Bind(ctx)

			suite.Equal(actualErr, tc.expectedErr)
			suite.Equal(bodyValidator.Model, tc.expectedBindModel)
		}

		suite.Run(tc.description, func() {
			testutil.Response(
				suite.T(),
				suite.router,
				"POST",
				"/test",
				tc.reqPayload,
			)
		})
	}
}

func TestUpdateTodoValidatorBind(t *testing.T) {
	reqBody := map[string]interface{}{
		"a": 5,
	}

	r := gin.New()
	r.POST("test", func(ctx *gin.Context) {
		bodyValidator := todo.NewUpdateTodoValidator()
		if err := bodyValidator.Bind(ctx); err != nil {
			t.Error(err)
		}

		// if _, ok := m["a"]; !ok {
		// 	t.Error("Not Exists Key a")
		// }
	})

	gin.SetMode(gin.TestMode)

	testutil.Response(
		t,
		r,
		"POST",
		"/test",
		testutil.ReqBodyFromInterface(t, reqBody),
	)
}
