package database

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisAPI interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}
