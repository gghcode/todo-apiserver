package auth

// PasswordAuthenticator authenticate password
type PasswordAuthenticator interface {
	IsValidPassword(password string, hash []byte) bool
}
