package redis

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/go-redis/redis"
)

type redisConn struct {
	client *redis.Client
}

// NewConnection return new connection of redis
func NewConnection(cfg config.Configuration) (Connection, func()) {
	conn := redisConn{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
		}),
	}

	cleanupFunc := func() {
		conn.Close()
	}

	return &conn, cleanupFunc
}

func (conn *redisConn) Healthy() bool {
	_, err := conn.client.Ping().Result()
	return err == nil
}

func (conn *redisConn) Client() *redis.Client {
	return conn.client
}

func (conn *redisConn) Close() error {
	return conn.client.Close()
}
