package cache

import (
	"context"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := NewMemory(10*time.Second, time.Second)
	ctx := context.Background()

	_ = c.Set(ctx, "foo", "bar")
	var value string
	err := c.Scan(ctx, "foo", &value)
	if err != nil {
		t.Fatal(err)
	}
	if value != "bar" {
		t.Error("value should be bar")
	}

	time.Sleep(12 * time.Second)

	err = c.Scan(ctx, "foo", &value)
	if err == nil {
		t.Error("foo should not exist")
	}

	_ = c.Set(ctx, "foo", "bar")
	_ = c.Remove(ctx, "foo")
	err = c.Scan(ctx, "foo", &value)
	if err == nil {
		t.Error("foo should not exist")
	}
}
