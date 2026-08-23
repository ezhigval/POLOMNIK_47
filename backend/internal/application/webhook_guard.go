package application

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"strings"
	"time"

	"polomnik/internal/ports"
)

const webhookIdempotencyTTL = 24 * time.Hour

func SecretEqual(got, want string) bool {
	left := strings.TrimSpace(got)
	right := strings.TrimSpace(want)
	if right == "" {
		return left == ""
	}
	if len(left) != len(right) {
		// Compare against itself so timing does not leak the expected length.
		subtle.ConstantTimeCompare([]byte(right), []byte(right))
		return false
	}
	return subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}

type WebhookGuard struct {
	cache ports.CachePort
}

func NewWebhookGuard(cache ports.CachePort) *WebhookGuard {
	return &WebhookGuard{cache: cache}
}

func (g *WebhookGuard) AlreadyProcessed(ctx context.Context, source, id string) bool {
	if g == nil || g.cache == nil {
		return false
	}
	key := webhookIdempotencyKey(source, id)
	if key == "" {
		return false
	}
	_, err := g.cache.Get(ctx, key)
	return err == nil
}

func (g *WebhookGuard) Remember(ctx context.Context, source, id string) {
	if g == nil || g.cache == nil {
		return
	}
	key := webhookIdempotencyKey(source, id)
	if key == "" {
		return
	}
	_ = g.cache.Set(ctx, key, []byte("1"), webhookIdempotencyTTL)
}

func webhookIdempotencyKey(source, id string) string {
	source = strings.TrimSpace(source)
	id = strings.TrimSpace(id)
	if source == "" || id == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(source + ":" + id))
	return "webhook:" + hex.EncodeToString(sum[:])
}
