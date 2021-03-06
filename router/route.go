package router

import (
	"net/http"
)

// RouteAction func
type RouteAction func(*http.Request) (string, error)

// Route struct
type Route struct {
	pattern string
	method  string
	Action  RouteAction
}

// Match method and URL are correct
func (r *Route) Match(method string, url string) bool {
	return r.method == method && r.pattern == url
}
