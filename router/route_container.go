package router

import "net/http"

// RouteContainer struct
type RouteContainer struct {
	routes []Route
}

// Get method route append to server
func (r *RouteContainer) Get(pattern string, action RouteAction) {
	r.AddRoute(http.MethodGet, pattern, action)
}

// Post method route append to server
func (r *RouteContainer) Post(pattern string, action RouteAction) {
	r.AddRoute(http.MethodPost, pattern, action)
}

// Put method route append to server
func (r *RouteContainer) Put(pattern string, action RouteAction) {
	r.AddRoute(http.MethodPut, pattern, action)
}

// Delete method route append to server
func (r *RouteContainer) Delete(pattern string, action RouteAction) {
	r.AddRoute(http.MethodDelete, pattern, action)
}

// AddRoute to server
func (r *RouteContainer) AddRoute(method string, pattern string, action RouteAction) {
	r.routes = append(r.routes, Route{
		method:  method,
		pattern: pattern,
		Action:  action,
	})
}

// Match routes
func (r *RouteContainer) Match(method string, url string) (Route, bool) {
	for _, route := range r.routes {
		if route.Match(method, url) {
			return route, false
		}
	}

	return Route{}, true
}
