package redis

import (
	"github.com/go-redis/redis"
)

// Connection is redis connection
type Connection interface {
	Healthy() bool
	Client() *redis.Client

	Close() error
}
