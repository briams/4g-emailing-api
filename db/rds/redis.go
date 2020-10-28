package rds

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Rds Redis model
type Rds struct {
	Ctx    context.Context
	Client *redis.Client
}

// NewRds creates a new Rds model
func NewRds(ctx context.Context, client *redis.Client) *Rds {
	return &Rds{ctx, client}
}

// RdsPing makes a ping to Redis
func (r *Rds) RdsPing() (string, error) {
	return r.Client.Ping(r.Ctx).Result()
}

// InfoToRAM set a new value to redis
func (r *Rds) InfoToRAM(key string, value interface{}) error {
	return r.Client.Set(r.Ctx, key, value, 0).Err()
}

// InfoFromRAM get a value from redis
func (r *Rds) InfoFromRAM(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

// KeysFromRAM get keys from redis
func (r *Rds) KeysFromRAM(pattern string) []string {
	return r.Client.Keys(r.Ctx, pattern).Val()
}

// DelFromRAM delete values from redis
func (r *Rds) DelFromRAM(keys ...string) error {
	return r.Client.Del(r.Ctx, keys...).Err()
}
