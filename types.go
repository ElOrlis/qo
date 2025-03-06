package qo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Query struct {
	Pairs  map[string]string
	Values []string
}

func (q Query) compile(r *http.Request) {
	var rawQuery string

	if len(q.Pairs) > 0 {
		url := r.URL.Query()
		for k, v := range q.Pairs {
			if len(v) > 0 {
				url.Add(k, v)
			}
		}
		rawQuery = url.Encode()
	}

	if len(q.Values) > 0 {
		encoded := make([]string, len(q.Values))
		for i, v := range q.Values {
			encoded[i] = url.QueryEscape(v)
		}

		values := strings.Join(encoded, "&")

		if len(rawQuery) > 0 {
			rawQuery += "&" + values
		} else {
			rawQuery = values
		}
	}

	r.URL.RawQuery = rawQuery
}

type Param struct {
	Name, Value string
}

type Params []Param

func (p Params) formatUrl(r *http.Request) {
	if len(p) > 0 {
		path := r.URL.Path
		for _, v := range p {
			path = strings.Replace(path, fmt.Sprintf(":%s", v.Name), v.Value, -1)
		}
		r.URL.Path = path
	}
}
