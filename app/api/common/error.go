package common

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
