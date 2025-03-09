package qo

import "net/http"

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Request struct {
	Client *http.Client
	Req    *http.Request
	Retry  *RetryPolicy
}
