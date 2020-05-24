package db

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/go-redis/redis"
)

// RedisConnection godoc
type RedisConnection interface {
	Client() *redis.Client
	Close() error
}

type redisConn struct {
	client *redis.Client
}

// NewRedisConn return new connection of redis
func NewRedisConn(cfg config.Configuration) RedisConnection {
	conn := redisConn{
		client: redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Addr,
		}),
	}

	return &conn
}

func (conn *redisConn) Client() *redis.Client {
	return conn.client
}

func (conn *redisConn) Close() error {
	return conn.client.Close()
}
