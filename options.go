package qo

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type BackOff interface {
	NextBackOff() time.Duration
	Reset()
}

func RetryStatusCodes(s []int) func(*Options) error {
	return func(o *Options) error {
		o.sc = s
		return nil
	}
}

func WithBackOff(b BackOff) func(*Options) error {
	return func(o *Options) error {
		o.backOff = b
		return nil
	}
}

func Body(b []byte) func(*http.Client, *http.Request) error {
	return func(_ *http.Client, r *http.Request) error {
		if b != nil {
			r.Body = io.NopCloser(strings.NewReader(string(b)))
		}
		return nil
	}
}

func Header(h http.Header) func(*http.Client, *http.Request) error {
	return func(_ *http.Client, r *http.Request) error {
		if h != nil {
			r.Header = h
		}
		return nil
	}
}

func QueryValues(q Query) func(*http.Client, *http.Request) error {
	return func(_ *http.Client, r *http.Request) error {
		q.compile(r)
		return nil
	}
}

func UrlParams(p Params) func(*http.Client, *http.Request) error {
	return func(_ *http.Client, r *http.Request) error {
		p.formatUrl(r)
		return nil
	}
}
