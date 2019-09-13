package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorMetaData godoc
type ErrorMetaData struct {
	Message string `json:"message"`
}

// ErrorResponse godoc
type ErrorResponse struct {
	Errors []ErrorMetaData `json:"errors"`
}

// NewErrRes godoc
func NewErrRes(err ...error) ErrorResponse {
	result := ErrorResponse{Errors: []ErrorMetaData{}}

	for _, err := range err {
		result.Errors = append(result.Errors, ErrorMetaData{Message: err.Error()})
	}

	return result
}

// HandledError godoc
type HandledError struct {
	Status int
	Errors []error
}

// NewHandledError godoc
func NewHandledError(code int, err ...error) *HandledError {
	return &HandledError{
		Status: code,
		Errors: err,
	}
}

func (err HandledError) Error() string {
	var errStrings []string

	for _, err := range err.Errors {
		errStrings = append(errStrings, err.Error())
	}

	return strings.Join(errStrings, "\n")
}

// WriteErrorResponse godoc
func WriteErrorResponse(ctx *gin.Context, err error) {
	handledErr, ok := err.(*HandledError)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, NewErrRes(err))
		return
	}

	ctx.JSON(handledErr.Status, NewErrRes(handledErr.Errors...))
}

// AbortErrorResponse godoc
func AbortErrorResponse(ctx *gin.Context, err error) {
	handledErr, ok := err.(HandledError)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewErrRes(err))
		return
	}

	ctx.AbortWithStatusJSON(handledErr.Status, NewErrRes(handledErr.Errors...))
}
