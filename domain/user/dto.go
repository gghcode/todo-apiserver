package user

import "time"

// CreateUserRequest is dto that contains info that require to create user.
type CreateUserRequest struct {
	UserName string `json:"username" validate:"required,min=4,max=100"`
	Password string `json:"password" validate:"required,min=8,max=50"`
}

// UserResponse is user response model.
type UserResponse struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"create_at"`
}
