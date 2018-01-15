package http

import (
	"net/http"
	"strconv"
)

type Server struct{}

func (s *Server) Listen(port int) (*http.Server) {
	srv := &http.Server{Addr: ":" + strconv.Itoa(port)}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Success"))
	})

	go func() {
		srv.ListenAndServe()
	}()

	return srv
}
