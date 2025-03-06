package qo

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cenkalti/backoff/v5"
)

func New(args ...func(*Client) error) (*Client, error) {
	return nil, nil
}

type Client struct {
	client           *http.Client
	backOffConfig    *backoff.ExponentialBackOff
	retryStatusCodes map[int]bool
}

func (c *Client) do(
	ctx context.Context,
	method, url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	op := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, backoff.Permanent(err)
		}

		for _, opt := range opts {
			err = opt(req)
			if err != nil {
				return nil, backoff.Permanent(err)
			}
		}

		var resp *http.Response
		resp, err = c.client.Do(req)
		if err != nil {
			return nil, backoff.Permanent(err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			var seconds int64
			seconds, err = strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
			if err != nil {
				return nil, backoff.Permanent(
					fmt.Errorf(
						"failed to parse Retry-After header into a valid integer, error: %s",
						err,
					),
				)
			}
			return nil, backoff.RetryAfter(int(seconds))
		}

		if c.retryStatusCodes[resp.StatusCode] {
			return nil, fmt.Errorf("retrying request, response status code: %d", resp.StatusCode)
		}

		return resp, nil
	}

	resp, err := backoff.Retry(
		ctx,
		op,
		backoff.WithBackOff(c.backOffConfig),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GET(
	ctx context.Context,
	url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, url, opts...)
}

func (c *Client) POST(
	ctx context.Context,
	url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, url, opts...)
}

func (c *Client) PUT(
	ctx context.Context,
	url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, url, opts...)
}

func (c *Client) PATCH(
	ctx context.Context,
	url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPatch, url, opts...)
}

func (c *Client) DELETE(
	ctx context.Context,
	url string,
	opts ...func(*http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, url, opts...)
}
