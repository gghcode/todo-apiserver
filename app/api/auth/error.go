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
		errors.New("Invalid user credential..."),
	)
)
