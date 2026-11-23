package cache

import (
	"sync"
	"time"
)

type Cache[T any] interface {
	// Get returns a cached response for the given key if present and not expired.
	Get(key string) (*Item[T], bool)

	// SetTTL caches the response with a specified TTL.
	SetTTL(key string, item Item[T], ttl time.Duration) error
}

type Item[T any] struct {
	Expiry time.Time
	Data   T
}

type InMem[T any] struct {
	data map[string]*Item[T]
	mu   sync.RWMutex
}

func (c *InMem[T]) SetTTL(key string, data T, ttl time.Duration) (err error) {
	item := &Item[T]{
		Expiry: time.Now().Add(ttl),
		Data:   data,
	}

	c.mu.Lock()
	c.data[key] = item
	c.mu.Unlock()

	return nil
}

func (c *InMem[T]) Get(key string) (*Item[T], bool) {
	c.mu.RLock()
	entry, found := c.data[key]
	c.mu.RUnlock()
	if !found {
		return nil, false
	}

	if time.Now().After(entry.Expiry) {
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return nil, false
	}

	return entry, true
}
