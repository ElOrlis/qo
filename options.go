package qo

import (
	"io"
	"net/http"
	"strings"
)

type CacheKeyFunc func(*http.Request) string

func WithCacheKeyFunc(f CacheKeyFunc) ClientOption {
	return func(c *Client) {
		c.cacheKeyFunc = f
	}
}

func defaultCacheKeyFunc(r *http.Request) string {
	return r.URL.String()
}

type ClientOption func(*Client)

func Body(b []byte) func(tx *Transaction) error {
	return func(tx *Transaction) error {
		if b != nil {
			tx.Request.Body = io.NopCloser(strings.NewReader(string(b)))
		}
		return nil
	}
}

func Header(h http.Header) func(tx *Transaction) error {
	return func(tx *Transaction) error {
		if h != nil {
			tx.Request.Header = h
		}
		return nil
	}
}

func QueryValues(q Query) func(tx *Transaction) error {
	return func(tx *Transaction) error {
		q.compile(tx.Request)
		return nil
	}
}

func UrlParams(p Params) func(tx *Transaction) error {
	return func(tx *Transaction) error {
		p.formatUrl(tx.Request)
		return nil
	}
}
