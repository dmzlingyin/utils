package memory

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(10*time.Second, time.Second)
	c.Set("foo", "bar")
	v, ok := c.Get("foo")
	if !ok {
		t.Error("foo should exist")
	}
	if v != "bar" {
		t.Error("value should be bar")
	}
	time.Sleep(12 * time.Second)
	v, ok = c.Get("foo")
	if ok {
		t.Error("foo should not exist")
	}

	c.Set("foo", "bar")
	c.Remove("foo")
	v, ok = c.Get("foo")
	if ok {
		t.Error("foo should not exist")
	}
}
