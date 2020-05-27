package db

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/go-redis/redis"
)

// RedisConnection godoc
type RedisConnection interface {
	Healthy() bool
	Client() *redis.Client

	Close() error
}

type redisConn struct {
	client *redis.Client
}

// NewRedisConn return new connection of redis
func NewRedisConn(cfg config.Configuration) (RedisConnection, func()) {
	conn := redisConn{
		client: redis.NewClient(&redis.Options{
			Addr: cfg.Redis.Addr,
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
