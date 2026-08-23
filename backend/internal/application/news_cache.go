package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/google/uuid"

	"polomnik/internal/ports"
)

const newsNamespaceKey = "news:namespace"

func (s *NewsService) invalidateNewsCache(ctx context.Context) {
	if s.cache == nil {
		return
	}
	_ = s.cache.Set(ctx, newsNamespaceKey, []byte(uuid.NewString()), defaultCacheTTL)
}

func (s *NewsService) cacheNamespace(ctx context.Context) (string, error) {
	value, err := s.cache.Get(ctx, newsNamespaceKey)
	if err == nil {
		return string(value), nil
	}
	if err != nil && !errors.Is(err, ports.ErrCacheMiss) {
		return "", err
	}
	namespace := uuid.NewString()
	_ = s.cache.Set(ctx, newsNamespaceKey, []byte(namespace), defaultCacheTTL)
	return namespace, nil
}

func (s *NewsService) cacheKey(ctx context.Context, parts ...string) (string, error) {
	namespace, err := s.cacheNamespace(ctx)
	if err != nil {
		return "", err
	}
	raw := namespace
	for _, part := range parts {
		raw += ":" + part
	}
	sum := sha256.Sum256([]byte(raw))
	return "news:" + hex.EncodeToString(sum[:]), nil
}

func newsListCacheID(pagination ports.Pagination) string {
	payload, _ := json.Marshal(pagination)
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

func (s *NewsService) readCache(ctx context.Context, key string, target any) (bool, error) {
	value, err := s.cache.Get(ctx, key)
	if err != nil {
		return false, nil
	}
	if err := json.Unmarshal(value, target); err != nil {
		return false, nil
	}
	return true, nil
}

func (s *NewsService) writeCache(ctx context.Context, key string, value any) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.cache.Set(ctx, key, payload, defaultCacheTTL)
}
