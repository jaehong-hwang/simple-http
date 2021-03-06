package http

import (
	"net/http"
	"strconv"

	"github.com/jaehong-hwang/simple-http/router"
)

// Server struct
type Server struct {
	Router router.RouteContainer
	Env    *ServerEnv
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, notFound := s.Router.Match(req.Method, req.URL.Path)

	if notFound == true {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	contents, err := route.Action(req)

	if err != nil {
		if serverErr, ok := err.(ServerError); ok {
			w.WriteHeader(serverErr.Status)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	w.Write([]byte(contents))
}

// Listen start
func (s *Server) Listen() *http.ServeMux {
	srv := http.NewServeMux()
	srv.HandleFunc("/", s.ServeHTTP)

	http.ListenAndServe(":"+strconv.Itoa(s.Env.Port), srv)

	return srv
}
