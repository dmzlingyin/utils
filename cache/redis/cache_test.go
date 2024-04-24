package redis

import (
	"context"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	url := "redis://:@192.168.7.251:6379/0"
	c := New(url, 1*time.Minute)

	ctx := context.Background()
	if err := c.Set(ctx, "foo", "bar"); err != nil {
		t.Fatal(err)
	}

	var v string
	if err := c.Scan(ctx, "foo", &v); err != nil {
		t.Fatal(err)
	}
	if v != "bar" {
		t.Fatal("value not match")
	}

	exists, err := c.Exists(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("key not exists")
	}
}
