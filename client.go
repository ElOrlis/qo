package qo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cenkalti/backoff/v5"
)

func New() (*Client, error) {
	return nil, nil
}

type Client struct {
	client *http.Client
	retry  *RetryPolicy
	logr   Logger
}

func defaultRetryHook(r *Transaction) func() (*http.Response, error) {
	return func() (*http.Response, error) {
		resp, err := r.Client.Do(r.Request)
		if err != nil {
			return nil, backoff.Permanent(err)
		}
		return resp, nil
	}
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
	tx := Transaction{Client: c.client, Request: req, Retry: c.retry}

	for _, opt := range opts {
		err := opt(&tx)
		if err != nil {
			return nil, err
		}
	}

	var retryHook backoff.Operation[*http.Response]

	switch tx.Retry.Hook != nil {
	case true:
		retryHook = defaultRetryHook(&tx)
	case false:
		retryHook = tx.Retry.Hook(tx.Client, tx.Request)
	}

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
		retryHook,
		args...,
	)
	if err != nil {
		return nil, err
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
