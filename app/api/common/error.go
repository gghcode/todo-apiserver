package common

// APIHandledError godoc
type APIHandledError struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

// APIHandledError2 godoc
type APIHandledError2 struct {
	Status int
	Errors []error
}

func (err *APIHandledError2) Error() string {
	return err.Errors[0].Error()
}

// Response godoc
func (err *APIHandledError2) Response() ErrorResponse {
	errs := []APIError{}

	for _, err := range err.Errors {
		errs = append(errs, APIError{
			Message: err.Error(),
		})
	}

	return ErrorResponse{
		Errors: errs,
	}
}

func (err *APIHandledError) Error() string {
	return err.Message
}

// Response godoc
func (err *APIHandledError) Response() ErrorResponse {
	return ErrorResponse{
		Errors: []APIError{
			{
				Message: err.Message,
			},
		},
	}
}

// APIError is http error object.
type APIError struct {
	Message string `json:"message"`
}

// NewErrResp is return new error response.
func NewErrResp(err error) ErrorResponse {
	return ErrorResponse{
		Errors: []APIError{
			{
				Message: err.Error(),
			},
		},
	}
}
