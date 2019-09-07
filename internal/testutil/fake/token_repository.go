package fake

import (
	"time"

	"github.com/stretchr/testify/mock"
)

// TokenRepository godoc
type TokenRepository struct {
	mock.Mock
}

// SaveRefreshToken godoc
func (repo *TokenRepository) SaveRefreshToken(userID int64, token string, expireIn time.Duration) error {
	args := repo.Called(userID, token, expireIn)
	return args.Error(0)
}

// UserIDByRefreshToken godoc
func (repo *TokenRepository) UserIDByRefreshToken(t string) (int64, error) {
	args := repo.Called(t)
	return args.Get(0).(int64), args.Error(1)
}
