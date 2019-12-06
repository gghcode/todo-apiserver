package auth

import (
	"errors"
)

var (
	// ErrInvalidCredential godoc
	ErrInvalidCredential = errors.New("Invalid user credential")

	// ErrNotContainToken godoc
	ErrNotContainToken = errors.New("Not contain token")

	// ErrInvalidToken godoc
	ErrInvalidToken = errors.New("Invalid token")

	// ErrInvalidTokenType godoc
	ErrInvalidTokenType = errors.New("Invalid token type")

	// ErrTokenExpired godoc
	ErrTokenExpired = errors.New("Token is expired")

	// ErrNotStoredToken godoc
	ErrNotStoredToken = errors.New("Token is not stored")
)