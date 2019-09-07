package auth

import (
	"errors"
	"net/http"

	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

var (
	// ErrInvalidCredential godoc
	ErrInvalidCredential = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Invalid user credential"),
	)

	// ErrNotContainToken godoc
	ErrNotContainToken = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Not contain token"),
	)

	// ErrInvalidToken godoc
	ErrInvalidToken = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Invalid token"),
	)

	// ErrInvalidTokenType godoc
	ErrInvalidTokenType = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Invalid token type"),
	)

	// ErrTokenExpired godoc
	ErrTokenExpired = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Token is expired"),
	)

	// ErrNotStoredToken godoc
	ErrNotStoredToken = api.NewHandledError(
		http.StatusUnauthorized,
		errors.New("Token is not stored"),
	)
)
