package todo_test

import (
	"io"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/api/todo"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
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
	wrongContents := "1"
	fakeModel := todo.UpdateTodoRequest{
		Contents: &wrongContents,
	}

	testCases := []struct {
		description       string
		reqPayload        io.Reader
		expectedErr       error
		expectedBindModel todo.UpdateTodoRequest
	}{
		{
			description: "ShouldBindModelWithTitle",
			reqPayload: testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
				"title": &fakeTitle,
			}),
			expectedErr: nil,
			expectedBindModel: todo.UpdateTodoRequest{
				Title: &fakeTitle,
			},
		},
		{
			description: "ShouldReturnErrWhenInvalidFieldType",
			reqPayload: testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
				"contents": 5,
			}),
			expectedErr: todo.JsonTypeError{
				Value: "number",
				Field: "contents",
			},
		},
		{
			description: "ShouldReturnErrWhenNotContentsMin",
			reqPayload: testutil.ReqBodyFromInterface(suite.T(), map[string]interface{}{
				"contents": "1",
			}),
			expectedErr: validation.ValidateStruct(&fakeModel,
				validation.Field(&fakeModel.Contents, validation.Length(5, 10)),
			),
			expectedBindModel: todo.UpdateTodoRequest{},
		},
	}

	for _, tc := range testCases {
		testFunc = func(ctx *gin.Context) {
			bodyValidator := todo.NewUpdateTodoValidator()
			actualErr := bodyValidator.Bind(ctx)

			suite.Equal(tc.expectedErr, actualErr)

			if tc.expectedErr == nil {
				suite.Equal(tc.expectedBindModel, bodyValidator.Model)
			}
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
