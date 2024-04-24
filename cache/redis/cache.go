package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	client *redis.Client
	ttl    time.Duration
}

func New(url string, ttl time.Duration) *Cache {
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return &Cache{
		client: redis.NewClient(opt),
		ttl:    ttl,
	}
}

func (c *Cache) Set(ctx context.Context, key string, value any) error {
	return c.SetWithTTL(ctx, key, value, c.ttl)
}

func (c *Cache) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, v, ttl).Err()
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count == 1, err
}

func (c *Cache) Remove(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *Cache) Scan(ctx context.Context, key string, value any) error {
	v, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(v, value)
}
