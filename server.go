package http

import (
	"net/http"
	"strconv"
)

// Server struct
type Server struct {
	routes []Route
}

// Get method route append to server
func (s *Server) Get(pattern string, action RouteAction) {
	s.AddRoute(http.MethodGet, pattern, action)
}

// Post method route append to server
func (s *Server) Post(pattern string, action RouteAction) {
	s.AddRoute(http.MethodPost, pattern, action)
}

// Put method route append to server
func (s *Server) Put(pattern string, action RouteAction) {
	s.AddRoute(http.MethodPut, pattern, action)
}

// Delete method route append to server
func (s *Server) Delete(pattern string, action RouteAction) {
	s.AddRoute(http.MethodDelete, pattern, action)
}

// AddRoute to server
func (s *Server) AddRoute(method string, pattern string, action RouteAction) {
	s.routes = append(s.routes, Route{
		method:  method,
		pattern: pattern,
		Action:  action,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	notFound := true
	for _, route := range s.routes {
		if route.Match(req.Method, req.URL.Path) {
			contents, err := route.Action(req)

			if err != nil {
				if serverErr, ok := err.(ServerError); ok {
					w.WriteHeader(serverErr.Status)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}

			w.Write([]byte(contents))

			notFound = false
		}
	}

	if notFound == true {
		w.WriteHeader(http.StatusNotFound)
	}
}

// Listen start
func (s *Server) Listen(port int) *http.ServeMux {
	srv := http.NewServeMux()
	srv.HandleFunc("/", s.ServeHTTP)

	http.ListenAndServe(":"+strconv.Itoa(port), srv)

	return srv
}
