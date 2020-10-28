package rds

import (
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	once        sync.Once
	redisClient *redis.Client
)

// NewClient creates a new client for Redis
func NewClient(opts *redis.Options) *redis.Client {
	once.Do(func() {
		redisClient = redis.NewClient(opts)
	})

	return redisClient
}
