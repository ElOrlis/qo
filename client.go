package qo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/cenkalti/backoff/v5"
)

// func New(args ...func(*Client) error) (*Client, error) {
// 	c := Client{
// 		client:        &http.Client{},
// 		backOffConfig: backoff.NewExponentialBackOff(),
// 	}
// 	for _, arg := range args {
// 		err := arg(&c)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return &c, nil
// }

type Client struct {
	c   *http.Client
	bck backoff.BackOff
	sc  []int
}

func (c *Client) do(
	ctx context.Context,
	method, url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	op := func() (*http.Response, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, nil)
		if err != nil {
			return nil, backoff.Permanent(err)
		}

		for _, opt := range opts {
			err = opt(c.c, req)
			if err != nil {
				return nil, backoff.Permanent(err)
			}
		}

		var resp *http.Response
		resp, err = c.c.Do(req)
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

		var b []byte
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				return
			}
		}()

		if !slices.Contains(c.sc, resp.StatusCode) {
			return nil, fmt.Errorf(
				"retrying, response status code: %d, body: %s",
				resp.StatusCode, string(b),
			)
		}

		resp.Body = io.NopCloser(strings.NewReader(string(b)))

		return resp, nil
	}

	resp, err := backoff.Retry(
		ctx,
		op,
		backoff.WithBackOff(c.bck),
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *Client) GET(
	ctx context.Context,
	url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodGet, url, opts...)
}

func (c *Client) POST(
	ctx context.Context,
	url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPost, url, opts...)
}

func (c *Client) PUT(
	ctx context.Context,
	url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPut, url, opts...)
}

func (c *Client) PATCH(
	ctx context.Context,
	url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodPatch, url, opts...)
}

func (c *Client) DELETE(
	ctx context.Context,
	url string,
	opts ...func(*http.Client, *http.Request) error,
) (*http.Response, error) {
	return c.do(ctx, http.MethodDelete, url, opts...)
}
