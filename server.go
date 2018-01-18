package http

import (
	"net/http"
	"strconv"
)

type Server struct {
	routes []Route
}

func (s *Server) AddRoute(method string, pattern string, action RouteAction) {
	s.routes = append(s.routes, Route{
		method: method,
		pattern: pattern,
		Action: action,
	})
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	notFound := true
	for _, route := range s.routes {
		if route.Match(req.Method, req.URL.Path) {
			res.Write(route.Action(req))
			notFound = false
		}
	}

	if notFound == true {
		res.WriteHeader(http.StatusNotFound)
	}
}

func (s *Server) Listen(port int) (*http.ServeMux) {
	srv := http.NewServeMux()
	srv.HandleFunc("/", s.ServeHTTP)

	http.ListenAndServe(":"+strconv.Itoa(port), srv)

	return srv
}
