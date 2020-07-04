package security

// Passport is object that execute about password auth
type Passport interface {
	HashPassword(password string) ([]byte, error)
	IsValidPassword(password string, hash []byte) bool
}
