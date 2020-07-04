package repository

import (
	"errors"
	"strconv"
	"time"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	"github.com/go-redis/redis"
)

var (
	// ErrRedisConnection is redis command error
	ErrRedisConnection = errors.New("redis command was failed")
)

type redisTokenRepository struct {
	redisConn db.RedisConnection
}

// NewRedisTokenRepository godoc
func NewRedisTokenRepository(redisConn db.RedisConnection) auth.TokenRepository {
	return &redisTokenRepository{
		redisConn: redisConn,
	}
}

func (repo *redisTokenRepository) SaveRefreshToken(userID int64, token string, expireIn time.Duration) error {
	return repo.redisConn.Client().Set(
		redisRefreshTokenKey(token),
		userID,
		expireIn*time.Second,
	).Err()
}

func (repo *redisTokenRepository) UserIDByRefreshToken(refreshToken string) (int64, error) {
	userIDStr, err := repo.redisConn.
		Client().
		Get(redisRefreshTokenKey(refreshToken)).
		Result()

	if err == redis.Nil {
		return -1, auth.ErrNotStoredToken
	} else if err != nil {
		return -1, ErrRedisConnection
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)

	return userID, err
}

func redisRefreshTokenKey(token string) string {
	const prefixRefreshToken = "refresh_token"
	return prefixRefreshToken + "-" + token
}
