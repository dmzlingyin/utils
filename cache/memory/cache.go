package memory

import (
	"runtime"
	"sync"
	"time"
)

const (
	DefaultTTL  = 10 * time.Minute
	InfiniteTTL = -1
)

type Cache struct {
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

func New(ttl time.Duration, evictInterval time.Duration) *Cache {
	c := &cache{
		elements: make(map[string]*Element),
		ttl:      ttl,
	}
	if ttl <= 0 {
		c.ttl = InfiniteTTL
	}

	C := &Cache{c}
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

func (c *cache) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elements[key] = &Element{
		value:  value,
		expiry: time.Now().Add(c.ttl).UnixNano(),
	}
}

func (c *cache) SetWithTTL(key string, value any, ttl time.Duration) {
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

func (c *cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.elements[key]
	return ok
}

func (c *cache) Remove(key string) {
	if !c.Exists(key) {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.elements, key)
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

func shutdown(c *Cache) {
	c.evictor.stop <- true
}
