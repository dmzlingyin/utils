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
	_, ok := c.elements[key]
	return ok, nil
}

func (c *cache) Remove(ctx context.Context, key string) error {
	ok, err := c.Exists(ctx, key)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.elements, key)
	return nil
}

func (c *cache) Scan(_ context.Context, key string, value any) error {
	v, ok := c.Get(key)
	if !ok {
		return errors.New("key not found in cache")
	}

	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		// 确保value是一个非空指针
		return errors.New("value must be a non-nil pointer")
	}
	valElem := val.Elem()
	if valElem.Type() != reflect.TypeOf(v) {
		// 确保缓存值类型与value指针指向的类型相同
		return errors.New("type mismatch")
	}
	valElem.Set(reflect.ValueOf(v))

	return nil
}

func (c *cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ele, ok := c.elements[key]
	if !ok {
		return nil, false
	}
	if ele.expiry > 0 && time.Now().UnixNano() > ele.expiry {
		return nil, false
	}
	return ele.value, true
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
