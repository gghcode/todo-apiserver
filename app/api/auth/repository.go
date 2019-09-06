package auth

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
)

const prefixRefreshToken = "refresh_token"

// Repository godoc
type Repository interface {
	SaveRefreshToken(userID int64, token string, expireIn time.Duration) error

	RefreshToken(userID int64) (string, error)
}

type repository struct {
	redisConn db.RedisConn
}

// NewRepository godoc
func NewRepository(redisConn db.RedisConn) Repository {
	return &repository{
		redisConn: redisConn,
	}
}

func (repo *repository) SaveRefreshToken(userID int64, token string, expireIn time.Duration) error {
	err := repo.redisConn.Client().Set(
		redisRefreshTokenKey(userID),
		token,
		expireIn*time.Second,
	).Err()

	return err
}

func (repo *repository) RefreshToken(userID int64) (string, error) {
	token, err := repo.redisConn.
		Client().
		Get(redisRefreshTokenKey(userID)).
		Result()

	if err == redis.Nil {
		return "", ErrNotStoredToken
	}

	return token, err
}

func redisRefreshTokenKey(userID int64) string {
	return fmt.Sprintf("%s_%d", prefixRefreshToken, userID)
}
