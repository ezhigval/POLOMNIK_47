package memory

import (
	"context"
	"testing"
	"time"
)

func TestLimiterBlocksAfterMax(t *testing.T) {
	limiter := New()
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		result, err := limiter.Allow(ctx, "ip", 2, time.Minute)
		if err != nil || !result.Allowed {
			t.Fatalf("request %d should be allowed: %+v %v", i+1, result, err)
		}
	}
	result, err := limiter.Allow(ctx, "ip", 2, time.Minute)
	if err != nil {
		t.Fatalf("allow: %v", err)
	}
	if result.Allowed {
		t.Fatal("expected block")
	}
	if result.RetryAfter <= 0 {
		t.Fatal("expected retry-after")
	}
}
