package auth

import (
	"errors"
)

var (
	// ErrNotContainTokenInHeader is error
	ErrNotContainTokenInHeader = errors.New("Token doesn't contain in header")

	// ErrInvalidCredential godoc
	ErrInvalidCredential = errors.New("Invalid user credential")

	// ErrInvalidToken godoc
	ErrInvalidToken = errors.New("Invalid token")

	// ErrInvalidTokenType godoc
	ErrInvalidTokenType = errors.New("Invalid token type")

	// ErrTokenExpired godoc
	ErrTokenExpired = errors.New("Token is expired")

	// ErrNotStoredToken godoc
	ErrNotStoredToken = errors.New("Token is not stored")
)
