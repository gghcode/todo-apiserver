package db

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/go-redis/redis"
)

// RedisConn godoc
type RedisConn interface {
	Client() *redis.Client
}

type redisConn struct {
	client *redis.Client
}

func (conn *redisConn) Client() *redis.Client {
	return conn.client
}

// NewRedisConn return new connection of redis
func NewRedisConn(cfg config.Configuration) RedisConn {
	conn := redisConn{
		client: redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Addr,
		}),
	}

	return &conn
}
