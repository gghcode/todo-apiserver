package user

import (
	"errors"
	"net/http"

	"gitlab.com/gyuhwan/apas-todo-apiserver/app/api"
)

var (
	// ErrAlreadyExistUser godoc
	ErrAlreadyExistUser = errors.New("Already exists user")

	// ErrUserNotFound godoc
	ErrUserNotFound = errors.New("User was not found")

	// ErrInvalidUserID godoc
	ErrInvalidUserID = api.NewHandledError(
		http.StatusBadRequest,
		errors.New("Invalid UserID"),
	)
	// .APIHandledError{
	// 	StatusCode: http.StatusBadRequest,
	// 	Message:    "Invalid UserID",
	// }
)
