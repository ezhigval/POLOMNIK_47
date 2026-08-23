package memory

import (
	"context"
	"sync"
	"time"

	"palomnik/internal/ports"
)

type counter struct {
	count    int
	windowAt time.Time
}

type Limiter struct {
	mu      sync.Mutex
	entries map[string]counter
}

func New() *Limiter {
	return &Limiter{entries: make(map[string]counter)}
}

var _ ports.RateLimiter = (*Limiter)(nil)

func (l *Limiter) Allow(_ context.Context, key string, max int, window time.Duration) (ports.RateLimitResult, error) {
	if max < 1 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	item, ok := l.entries[key]
	if !ok || now.Sub(item.windowAt) >= window {
		l.entries[key] = counter{count: 1, windowAt: now}
		return ports.RateLimitResult{Allowed: true}, nil
	}

	item.count++
	l.entries[key] = item
	if item.count > max {
		retry := window - now.Sub(item.windowAt)
		if retry < time.Second {
			retry = time.Second
		}
		return ports.RateLimitResult{Allowed: false, RetryAfter: retry}, nil
	}
	return ports.RateLimitResult{Allowed: true}, nil
}
