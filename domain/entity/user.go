package entity

// User is user entity mode
type User struct {
	ID           int64
	UserName     string
	PasswordHash []byte
}
