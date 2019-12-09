package user

import (
	"errors"
)

var (
	// ErrAlreadyExistUser godoc
	ErrAlreadyExistUser = errors.New("Already exists user")

	// ErrUserNotFound godoc
	ErrUserNotFound = errors.New("User was not found")

	// ErrInvalidUserID godoc
	ErrInvalidUserID = errors.New("Invalid UserID")
)
