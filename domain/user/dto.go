package user

import "time"

// UserResponse is user response model.
type UserResponse struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"username"`
	CreatedAt time.Time `json:"create_at"`
}
