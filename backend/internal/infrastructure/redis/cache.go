package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (cache *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return cache.client.Get(ctx, key).Result()
}

func (cache *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return cache.client.Set(ctx, key, value, ttl).Err()
}
