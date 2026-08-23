package noop

import (
	"context"
	"time"

	"palomnik/internal/ports"
)

type Cache struct{}

func New() Cache {
	return Cache{}
}

func (Cache) Get(context.Context, string) ([]byte, error) {
	return nil, ports.ErrCacheMiss
}

func (Cache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (Cache) Delete(context.Context, string) error {
	return nil
}
