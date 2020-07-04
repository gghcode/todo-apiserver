package user

// UseCase is user usecase interface
type UseCase interface {
	CreateUser(CreateUserRequest) (UserResponse, error)
	GetUserByUserID(userID int64) (UserResponse, error)
	GetUserByUserName(userName string) (UserResponse, error)
}
