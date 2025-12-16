package cache

import (
	"Gobank/pkg/config"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	SetIdempotencyKey(ctx context.Context, key string, result string, ttl time.Duration) (bool, error)
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(cfg *config.Config) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return &redisCache{client: client}
}
func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}
func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
func (c *redisCache) SetIdempotencyKey(ctx context.Context, key string, result string, ttl time.Duration) (bool, error) {
	success, err := c.client.SetNX(ctx, "idemp:"+key, result, ttl).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}
