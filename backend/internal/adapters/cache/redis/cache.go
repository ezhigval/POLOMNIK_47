package rediscache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"palomnik/internal/ports"
)

type Cache struct {
	client *redis.Client
}

func New(redisURL string) (*Cache, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return &Cache{client: client}, nil
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ports.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) Ping(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("cache is not configured")
	}
	return c.client.Ping(ctx).Err()
}

func (c *Cache) Allow(ctx context.Context, key string, max int, window time.Duration) (ports.RateLimitResult, error) {
	if c == nil || c.client == nil {
		return ports.RateLimitResult{}, fmt.Errorf("cache is not configured")
	}
	if max < 1 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	redisKey := "rl:" + key
	count, err := c.client.Incr(ctx, redisKey).Result()
	if err != nil {
		return ports.RateLimitResult{}, err
	}
	if count == 1 {
		if err := c.client.Expire(ctx, redisKey, window).Err(); err != nil {
			return ports.RateLimitResult{}, err
		}
	}

	if count <= int64(max) {
		return ports.RateLimitResult{Allowed: true}, nil
	}

	ttl, err := c.client.TTL(ctx, redisKey).Result()
	if err != nil || ttl <= 0 {
		ttl = window
	}
	return ports.RateLimitResult{Allowed: false, RetryAfter: ttl}, nil
}
