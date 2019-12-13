package user

import "time"

// CreateUserRequest is dto that contains info that require to create user.
type CreateUserRequest struct {
	UserName string
	Password string
}

// UserResponse is user response model.
type UserResponse struct {
	ID        int64
	UserName  string
	CreatedAt time.Time
}
