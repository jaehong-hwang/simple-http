package http

import "net/http"

type RouteAction func(*http.Request) []byte

type Route struct {
	pattern string
	method string
	Action RouteAction
}

func (r *Route) Match(method string, url string) bool {
	return r.method == method && r.pattern == url
}
