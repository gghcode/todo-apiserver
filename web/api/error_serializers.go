package api

type (
	// ErrorMetaData godoc
	ErrorMetaData struct {
		Message string `json:"message"`
	}

	// ErrorResponse godoc
	ErrorResponse struct {
		Error ErrorMetaData `json:"error"`
	}
)

// MakeErrorResponse return error response from error
func MakeErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error: ErrorMetaData{
			Message: err.Error(),
		},
	}
}
