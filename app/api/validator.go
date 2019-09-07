package api

import (
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate godoc
func Validate(model interface{}) error {
	if err := validate.Struct(model); err != nil {
		return NewHandledError(http.StatusBadRequest, err)
	}

	return nil
}
