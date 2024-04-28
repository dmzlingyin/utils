package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedis(url string, ttl time.Duration) Cache {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return &Redis{
		client: redis.NewClient(opt),
		ttl:    ttl,
	}
}

func (r *Redis) Set(ctx context.Context, key string, value any) error {
	return r.SetWithTTL(ctx, key, value, r.ttl)
}

func (r *Redis) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, v, ttl).Err()
}

func (r *Redis) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count == 1, err
}

func (r *Redis) Remove(ctx context.Context, key string) error {
	exists, err := r.Exists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return ErrKeyNotFound
	}
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Scan(ctx context.Context, key string, value any) error {
	exists, err := r.Exists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return ErrKeyNotFound
	}
	v, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(v, value)
}
