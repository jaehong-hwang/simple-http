package http

import "net/http"

// Routes container
type Routes struct {
	routes []Route
}

// Get method route append to server
func (r *Routes) Get(pattern string, action RouteAction) {
	r.AddRoute(http.MethodGet, pattern, action)
}

// Post method route append to server
func (r *Routes) Post(pattern string, action RouteAction) {
	r.AddRoute(http.MethodPost, pattern, action)
}

// Put method route append to server
func (r *Routes) Put(pattern string, action RouteAction) {
	r.AddRoute(http.MethodPut, pattern, action)
}

// Delete method route append to server
func (r *Routes) Delete(pattern string, action RouteAction) {
	r.AddRoute(http.MethodDelete, pattern, action)
}

// AddRoute to server
func (r *Routes) AddRoute(method string, pattern string, action RouteAction) {
	r.routes = append(r.routes, Route{
		method:  method,
		pattern: pattern,
		Action:  action,
	})
}

// Match routes
func (r *Routes) Match(method string, url string) (Route, bool) {
	for _, route := range r.routes {
		if route.Match(method, url) {
			return route, false
		}
	}

	return Route{}, true
}
