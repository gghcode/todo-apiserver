package user

import (
	"errors"
	"net/http"

	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

var (
	// ErrAlreadyExistUser godoc
	ErrAlreadyExistUser = api.NewHandledError(
		http.StatusConflict,
		errors.New("Already exists user"),
	)

	// ErrUserNotFound godoc
	ErrUserNotFound = api.NewHandledError(
		http.StatusNotFound,
		errors.New("User was not found"),
	)

	// ErrInvalidUserID godoc
	ErrInvalidUserID = api.NewHandledError(
		http.StatusBadRequest,
		errors.New("Invalid UserID"),
	)
)
