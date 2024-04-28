package cache

import (
	"context"
	"testing"
	"time"
)

type User struct {
	Name string
}

func TestRedis(t *testing.T) {
	url := "redis://:@192.168.7.251:6379/0"
	c := NewRedis(url, 5*time.Second)

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

	if err = c.Remove(ctx, "foo"); err != nil {
		t.Fatal(err)
	}
	exists, err = c.Exists(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("key should not exists")
	}

	u := User{Name: "lingyin"}
	if err = c.Set(ctx, "user", u); err != nil {
		t.Fatal(err)
	}
	var u2 User
	if err = c.Scan(ctx, "user", &u2); err != nil {
		t.Fatal(err)
	}

	time.Sleep(6 * time.Second)
	exists, err = c.Exists(ctx, "foo")
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("key should not exists")
	}
}
