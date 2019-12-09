package api

import (
	"fmt"
	"net/http"
)

type (
	// HandledError is handled error
	HandledError interface {
		error
		StatusCode() int
	}

	// UnmarshalError is error
	UnmarshalError struct {
		Value string
		Field string
	}
)

// NewUnmarshalError return unmarshal error
func NewUnmarshalError(field, value string) error {
	return UnmarshalError{
		Field: field,
		Value: value,
	}
}

// StatusCode godoc
func (err UnmarshalError) StatusCode() int {
	return http.StatusBadRequest
}

func (err UnmarshalError) Error() string {
	return fmt.Sprintf("json: '%s' field is %s type", err.Field, err.Value)
}
