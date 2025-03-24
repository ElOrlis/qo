package qo

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

func WithCacheEnabled(enabled bool) ClientOption {
	return func(c *Client) {
		c.cacheEnabled = enabled
	}
}

func WithCache(cache Cache) ClientOption {
	return func(c *Client) {
		c.cache = cache
		c.cacheEnabled = true
	}
}

func WithCacheTTL(ttl time.Duration) ClientOption {
	return func(c *Client) {
		c.cacheTTL = ttl
	}
}

type CacheKeyFunc func(*http.Request) string

func WithCacheKeyFunc(f CacheKeyFunc) ClientOption {
	return func(c *Client) {
		c.cacheKeyFunc = f
	}
}

func defaultCacheKeyFunc(r *http.Request) string {
	return r.URL.String()
}

type Cache interface {
	// Get returns a cached response for the given key if present and not expired.
	Get(key string) (*http.Response, bool)

	// SetTTL caches the response with a specified TTL.
	SetTTL(key string, resp *http.Response, ttl time.Duration) error
}

type Item struct {
	Expiry   time.Time
	ReadBody []byte
	Response http.Response
}

func NewInMemoryCache() *InMemCache {
	return &InMemCache{
		data: make(map[string]*Item),
	}
}

type InMemCache struct {
	data map[string]*Item
	mu   sync.RWMutex
}

func (c *InMemCache) SetTTL(key string, resp *http.Response, ttl time.Duration) error {
	defer func(res *http.Response) {
		_ = res.Close
	}(resp)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = io.NopCloser(bytes.NewReader(b))

	item := &Item{
		Expiry:   time.Now().Add(ttl),
		Response: *resp,
		ReadBody: b,
	}

	c.mu.Lock()
	c.data[key] = item
	c.mu.Unlock()

	return nil
}

func (c *InMemCache) Get(key string) (*http.Response, bool) {
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

	resp := &entry.Response
	resp.Body = io.NopCloser(bytes.NewReader(entry.ReadBody))
	resp.Header = entry.Response.Header.Clone()

	return resp, true
}
