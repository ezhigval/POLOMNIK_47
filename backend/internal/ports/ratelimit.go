package ports

import (
	"context"
	"time"
)

type RateLimitResult struct {
	Allowed    bool
	RetryAfter time.Duration
}

// RateLimiter is optional. Redis down must fall back to an in-process limiter.
type RateLimiter interface {
	Allow(ctx context.Context, key string, max int, window time.Duration) (RateLimitResult, error)
}
