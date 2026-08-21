package fiberhttp

import "context"

type HealthDeps struct {
	PingDB        func(context.Context) error
	PingCache     func(context.Context) error
	CacheRequired bool
}
