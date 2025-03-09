package qo

import (
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v5"
)

var (
	Permanent  = backoff.Permanent
	RetryAfter = backoff.RetryAfter
)

type BackOff interface {
	NextBackOff() time.Duration
	Reset()
}

type HookOp func(HttpClient, *http.Request) func() (*http.Response, error)

type RetryPolicy struct {
	Hook   HookOp
	Policy BackOff
}
