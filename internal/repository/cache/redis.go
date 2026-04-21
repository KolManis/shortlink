package cache

import (
	"context"
	"log/slog"
	"time"

	urlDomain "github.com/KolManis/shortlink/internal/domain/url"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	logger *slog.Logger
}

func NewRedisCache(addr string, logger *slog.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{
		client: client,
		logger: logger,
	}, nil
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	c.logger.Debug("Redis GET", "key", key)

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.logger.Debug("Redis GET miss", "key", key)
		return "", urlDomain.ErrNotFound
	}
	if err != nil {
		c.logger.Error("Redis GET error", "key", key, "error", err)
		return "", err
	}

	c.logger.Debug("Redis GET hit", "key", key, "value", val)
	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	c.logger.Debug("Redis SET", "key", key, "ttl", ttl)
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	c.logger.Debug("Redis INCR", "key", key, "result", val)
	return val, err
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	c.logger.Debug("Redis Del")
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Close() error {
	c.logger.Debug("Redis Close")
	return c.client.Close()
}
