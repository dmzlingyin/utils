package cache

import (
	"context"
	"errors"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value any) error
	SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Remove(ctx context.Context, key string) error
	Scan(ctx context.Context, key string, value any) error
}

var (
	ErrKeyNotFound = errors.New("key not found in cache")
)
