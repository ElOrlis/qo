package qo

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
)

var (
	NewConstantBackoff    = backoff.NewConstantBackOff
	NewExponentialBackoff = backoff.NewExponentialBackOff

	Permanent  = backoff.Permanent
	RetryAfter = backoff.RetryAfter
)

func defaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		Hook: func(cli HttpClient, r *http.Request) func() (*http.Response, error) {
			return func() (*http.Response, error) {
				resp, err := cli.Do(r)
				if err != nil {
					return nil, Permanent(err)
				}
				return resp, nil
			}
		},
		MaxRetries:     3,
		MaxElapsedTime: 2 * time.Second,
		Policy:         NewExponentialBackoff(),
	}
}

type RetryPolicy struct {
	Hook           func(HttpClient, *http.Request) func() (*http.Response, error)
	MaxRetries     uint
	MaxElapsedTime time.Duration
	Notify         func(error, time.Duration)
	Policy         interface {
		NextBackOff() time.Duration
		Reset()
	}
}
