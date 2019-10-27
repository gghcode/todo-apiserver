package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	// ErrorMetaData godoc
	ErrorMetaData struct {
		Message string `json:"message"`
	}

	// ErrorResponse godoc
	ErrorResponse struct {
		Errors []ErrorMetaData `json:"errors"`
	}

	// handledError godoc
	handledError struct {
		Status int
		Errors []error
	}

	// HandledError godoc
	HandledError interface {
		StatusCode() int
		Error() string
	}

	// ErrorResponseFactory is factory that create error response.
	ErrorResponseFactory interface {
		ErrorResponse() ErrorResponse
	}
)

// NewErrRes godoc
func NewErrRes(err error) ErrorResponse {
	return ErrorResponse{
		Errors: []ErrorMetaData{
			{
				Message: err.Error(),
			},
		},
	}
}

// NewHandledError godoc
func NewHandledError(code int, err ...error) HandledError {
	return &handledError{
		Status: code,
		Errors: err,
	}
}

func (err handledError) StatusCode() int {
	return err.Status
}

func (err handledError) Error() string {
	var errStrings []string

	for _, err := range err.Errors {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, "\n")
}

// WriteErrorResponse godoc
func WriteErrorResponse(ctx *gin.Context, err error) {
	handledErr, ok := err.(HandledError)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, NewErrRes(err))
		return
	}

	if errResFactory, ok := err.(ErrorResponseFactory); ok {
		ctx.JSON(handledErr.StatusCode(), errResFactory.ErrorResponse())
		return
	}

	ctx.JSON(handledErr.StatusCode(), NewErrRes(handledErr))
}

// AbortErrorResponse godoc
func AbortErrorResponse(ctx *gin.Context, err error) {
	handledErr, ok := err.(HandledError)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrRes(err))
		return
	}

	ctx.AbortWithStatusJSON(handledErr.StatusCode(), NewErrRes(handledErr))
}
