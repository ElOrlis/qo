package qo

import (
	"context"
	"net/http"

	"github.com/cenkalti/backoff/v5"
)

type Client struct {
	client *http.Client
	retry  *RetryPolicy
}

func defaultRetryHook(r *Request) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		resp, err := r.Client.Do(r.Req)
		if err != nil {
			return nil, backoff.Permanent(err)
		}
		return resp, nil
	}
}

func (c *Client) do(
	ctx context.Context,
	method, url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	r := Request{Client: c.client, Req: req, Retry: c.retry}

	for _, opt := range opts {
		err := opt(&r)
		if err != nil {
			return nil, err
		}
	}

	var retryHook backoff.Operation[*http.Response]

	switch r.Retry.Hook != nil {
	case true:
		retryHook = defaultRetryHook(&r)
	case false:
		retryHook = r.Retry.Hook(r.Client, r.Req)
	}

	resp, err := backoff.Retry(ctx, retryHook, backoff.WithBackOff(r.Retry.Policy))
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GET(
	ctx context.Context,
	url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, url, opts...)
}

func (c *Client) POST(
	ctx context.Context,
	url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, url, opts...)
}

func (c *Client) PUT(
	ctx context.Context,
	url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, url, opts...)
}

func (c *Client) PATCH(
	ctx context.Context,
	url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPatch, url, opts...)
}

func (c *Client) DELETE(
	ctx context.Context,
	url string,
	opts ...func(*Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, url, opts...)
}
