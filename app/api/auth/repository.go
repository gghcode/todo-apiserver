package auth

import (
	"strconv"
	"time"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/go-redis/redis"
)

// Repository godoc
type Repository interface {
	SaveRefreshToken(userID int64, token string, expireIn time.Duration) error
	UserIDByRefreshToken(refreshToken string) (int64, error)
}

type repository struct {
	redisConn db.RedisConnection
}

// NewRepository godoc
func NewRepository(redisConn db.RedisConnection) Repository {
	return &repository{
		redisConn: redisConn,
	}
}

func (repo *repository) SaveRefreshToken(userID int64, token string, expireIn time.Duration) error {
	return repo.redisConn.Client().Set(
		redisRefreshTokenKey(token),
		userID,
		expireIn*time.Second,
	).Err()
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
	const prefixRefreshToken = "refresh_token"
	return prefixRefreshToken + "-" + token
}
