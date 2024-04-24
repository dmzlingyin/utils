package redis

import (
	"context"
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
	_, err := c.client.Set(ctx, key, value, c.ttl).Result()
	return err
}

func (c *Cache) SetWithTTL(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := c.client.Set(ctx, key, value, ttl).Result()
	return err
}

func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	return count == 1, err
}

func (c *Cache) Scan(ctx context.Context, key string, value any) error {
	return c.client.Get(ctx, key).Scan(value)
}
