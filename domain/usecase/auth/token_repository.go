package auth

import "time"

// TokenRepository godoc
type TokenRepository interface {
	SaveRefreshToken(userID int64, token string, expireIn time.Duration) error
	UserIDByRefreshToken(refreshToken string) (int64, error)
}
