package cache

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	InfiniteTTL = -1
)

type Memory struct {
	*cache
}

type cache struct {
	mu       sync.RWMutex
	elements map[string]*Element
	ttl      time.Duration
	evictor  *Evictor
}

type Element struct {
	value  any
	expiry int64
}

func NewMemory(ttl time.Duration, evictInterval time.Duration) Cache {
	c := &cache{
		elements: make(map[string]*Element),
		ttl:      ttl,
	}
	if ttl <= 0 {
		c.ttl = InfiniteTTL
	}

	C := &Memory{c}
	if evictInterval > 0 {
		c.evictor = &Evictor{
			interval: evictInterval,
			stop:     make(chan bool),
			cache:    c,
		}
		go c.evictor.run()
		runtime.SetFinalizer(C, shutdown)
	}

	return C
}

func (c *cache) Set(ctx context.Context, key string, value any) error {
	return c.SetWithTTL(ctx, key, value, c.ttl)
}

func (c *cache) SetWithTTL(_ context.Context, key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ele := &Element{
		value:  value,
		expiry: time.Now().Add(ttl).UnixNano(),
	}
	if ttl <= 0 {
		ele.expiry = InfiniteTTL
	}
	c.elements[key] = ele
	return nil
}

func (c *cache) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.elements[key]
	if !ok {
		return false, nil
	}
	if v.expiry > 0 && time.Now().UnixNano() > v.expiry {
		return false, nil
	}
	return true, nil
}

func (c *cache) Remove(ctx context.Context, key string) error {
	exists, err := c.Exists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.elements, key)
	return nil
}

func (c *cache) Scan(_ context.Context, key string, value any) error {
	v, ok := c.get(key)
	if !ok {
		return errors.New("key not found in cache")
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("value must be a non-nil pointer")
	}
	elem := val.Elem()
	if elem.Type() != reflect.TypeOf(v) {
		return errors.New("type mismatch")
	}
	elem.Set(reflect.ValueOf(v))

	return nil
}

func (c *cache) get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.elements[key]
	if !ok {
		return nil, false
	}
	if v.expiry > 0 && time.Now().UnixNano() > v.expiry {
		return nil, false
	}
	return v.value, true
}

type Evictor struct {
	interval time.Duration
	stop     chan bool
	cache    *cache
}

func (e *Evictor) run() {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.evict()
		case <-e.stop:
			return
		}
	}
}

func (e *Evictor) evict() {
	e.cache.mu.Lock()
	defer e.cache.mu.Unlock()
	for k, v := range e.cache.elements {
		if v.expiry > 0 && time.Now().UnixNano() > v.expiry {
			delete(e.cache.elements, k)
		}
	}
}

func shutdown(c *Memory) {
	c.evictor.stop <- true
}
