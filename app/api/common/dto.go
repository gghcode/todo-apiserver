package common

// ErrorResponse is app response.
type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}
