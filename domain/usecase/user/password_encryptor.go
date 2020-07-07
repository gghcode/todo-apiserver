package user

// PasswordEncryptor generate password hash
type PasswordEncryptor interface {
	HashPassword(password string) ([]byte, error)
}
