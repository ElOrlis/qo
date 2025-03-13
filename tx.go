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

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Transaction struct {
	Client  *http.Client
	Request *http.Request
	Retry   *RetryPolicy
}
