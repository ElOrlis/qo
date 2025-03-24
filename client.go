package qo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v5"
)

func New() (*Client, error) {
	return nil, nil
}

type Client struct {
	client       *http.Client
	retry        RetryPolicy
	logr         Logger
	cacheEnabled bool
	cacheTTL     time.Duration
	cache        Cache
	cacheKeyFunc CacheKeyFunc
}

func parseMaxAge(cc string) (time.Duration, error) {
	const prefix = "max-age="
	for _, part := range strings.Split(cc, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, prefix) {
			ttlStr := part[len(prefix):]
			ttl, err := time.ParseDuration(ttlStr + "s")
			if err != nil {
				return 0, err
			}
			return ttl, nil
		}
	}
	return 0, fmt.Errorf("max-age not found")
}

func (c *Client) do(
	ctx context.Context,
	method, url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	c.logr.Info(ctx, "Initiating request...")
	defer func() {
		c.logr.Info(ctx, "...request completed")
	}()
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		c.logr.Error(
			ctx,
			fmt.Sprintf("failed to create a request with context, error: %s", err),
			"error",
			err,
		)
		return nil, err
	}
	c.logr.Info(ctx, "successfully created http request")

	cli := *c.client
	retryPolicy := c.retry
	tx := Transaction{
		Client:       &cli,
		Request:      req,
		Retry:        &retryPolicy,
		CacheKeyFunc: c.cacheKeyFunc,
		CacheEnabled: c.cacheEnabled,
		CacheTTL:     c.cacheTTL,
	}

	for _, opt := range opts {
		err := opt(&tx)
		if err != nil {
			return nil, err
		}
	}

	if method == http.MethodGet && tx.CacheEnabled && c.cache != nil {
		var keyFunc CacheKeyFunc
		if tx.CacheKeyFunc != nil {
			keyFunc = tx.CacheKeyFunc
		} else {
			keyFunc = c.cacheKeyFunc
		}
		cacheKey := keyFunc(tx.Request)
		if resp, ok := c.cache.Get(cacheKey); ok {
			return resp, nil
		}
	}

	c.logr.Info(ctx, "created request transaction and added optional parameters")

	var args []backoff.RetryOption
	args = append(args,
		backoff.WithBackOff(tx.Retry.Policy),
		backoff.WithMaxTries(tx.Retry.MaxRetries),
		backoff.WithMaxElapsedTime(tx.Retry.MaxElapsedTime),
	)

	if tx.Retry.Notify != nil {
		args = append(args, backoff.WithNotify(tx.Retry.Notify))
	}

	resp, err := backoff.Retry(
		ctx,
		tx.Retry.Hook(tx.Client, tx.Request),
		args...,
	)
	if err != nil {
		c.logr.Error(ctx, "error while retrying operation", "Error", err)
		return nil, err
	}

	if method == http.MethodGet && tx.CacheEnabled && c.cache != nil {
		var keyFunc CacheKeyFunc
		if tx.CacheKeyFunc != nil {
			keyFunc = tx.CacheKeyFunc
		} else {
			keyFunc = c.cacheKeyFunc
		}
		cacheKey := keyFunc(tx.Request)
		ttl := tx.CacheTTL
		if cc := resp.Header.Get("Cache-Control"); cc != "" {
			if parsedTTL, err := parseMaxAge(cc); err == nil && parsedTTL > 0 {
				ttl = parsedTTL
			}
		}

		if err = c.cache.SetTTL(cacheKey, resp, ttl); err != nil {
			c.logr.Error(ctx, "failed to cache response", "Error", err)
		} else {
			c.logr.Info(ctx, "cached response for %s", cacheKey)
		}
	}

	return resp, nil
}

func (c *Client) GET(
	ctx context.Context,
	url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, url, opts...)
}

func (c *Client) POST(
	ctx context.Context,
	url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, url, opts...)
}

func (c *Client) PUT(
	ctx context.Context,
	url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, url, opts...)
}

func (c *Client) PATCH(
	ctx context.Context,
	url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPatch, url, opts...)
}

func (c *Client) DELETE(
	ctx context.Context,
	url string,
	opts ...func(*Transaction) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, url, opts...)
}
