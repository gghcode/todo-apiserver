package user

// User is user entity model
type User struct {
	ID           int64
	UserName     string
	PasswordHash []byte
	CreatedAt    int64
}
