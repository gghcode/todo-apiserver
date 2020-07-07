package repo

import (
	"errors"
	"strconv"
	"time"

	"github.com/gghcode/apas-todo-apiserver/domain/usecase/auth"
	myRedis "github.com/gghcode/apas-todo-apiserver/infra/redis"
	"github.com/go-redis/redis"
)

var (
	// ErrRedisConnection is redis command error
	ErrRedisConnection = errors.New("redis command was failed")
)

type redisTokenRepository struct {
	redisConn myRedis.Connection
}

// NewRedisTokenRepository godoc
func NewRedisTokenRepository(redisConn myRedis.Connection) auth.TokenRepository {
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
