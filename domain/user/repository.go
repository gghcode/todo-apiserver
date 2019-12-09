package user

// Repository godoc
type Repository interface {
	CreateUser(User) (User, error)
	AllUsers() ([]User, error)
	UserByID(userID int64) (User, error)
	UserByUserName(username string) (User, error)
	UpdateUserByID(user User) (User, error)
	RemoveUserByID(userID int64) (User, error)
}
