package user

// UsecaseInteractor is user usecase interface
type UsecaseInteractor interface {
	CreateUser(CreateUserRequest) (UserResponse, error)
	GetUserByUserID(userID int64) (UserResponse, error)
	GetUserByUserName(userName string) (UserResponse, error)
}
