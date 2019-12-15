package api

type (
	// ErrorMetaData godoc
	ErrorMetaData struct {
		Message string `json:"message"`
	}

	// ErrorResponseDTO godoc
	ErrorResponseDTO struct {
		Error ErrorMetaData `json:"error"`
	}
)

// MakeErrorResponseDTO return error response from error
func MakeErrorResponseDTO(err error) ErrorResponseDTO {
	return ErrorResponseDTO{
		Error: ErrorMetaData{
			Message: err.Error(),
		},
	}
}
