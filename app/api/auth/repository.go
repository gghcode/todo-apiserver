package auth

import (
	"strconv"
	"time"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/go-redis/redis"
)

const prefixRefreshToken = "refresh_token"

// Repository godoc
type Repository interface {
	SaveRefreshToken(userID int64, token string, expireIn time.Duration) error

	UserIDByRefreshToken(refreshToken string) (int64, error)
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
		redisRefreshTokenKey(token),
		userID,
		expireIn*time.Second,
	).Err()

	return err
}

func (repo *repository) UserIDByRefreshToken(refreshToken string) (int64, error) {
	userIDStr, err := repo.redisConn.
		Client().
		Get(redisRefreshTokenKey(refreshToken)).
		Result()

	if err == redis.Nil {
		return 0, ErrNotStoredToken
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	return userID, err
}

func redisRefreshTokenKey(token string) string {
	return prefixRefreshToken + "-" + token
}
