package memory

import (
	"runtime"
	"sync"
	"time"
)

const (
	DefaultTTL           = 10 * time.Minute
	DefaultCleanInterval = 10 * time.Second
	InfiniteTTL          = -1
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
	value interface{}
	ttl   time.Duration
}

func New(ttl time.Duration, evictInterval time.Duration) *Cache {
	c := &cache{
		elements: make(map[string]*Element),
		ttl:      ttl,
	}
	if ttl == 0 {
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
		runtime.SetFinalizer(C, c.evictor.shutdown)
	}

	return C
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
			e.clean()
		case <-e.stop:
			return
		}
	}
}

func (e *Evictor) clean() {

}

func (e *Evictor) shutdown() {
	e.stop <- true
}
