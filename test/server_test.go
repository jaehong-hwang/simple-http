package test

import (
	"testing"
	"flag"
	"net/http"
	simpleHttp "bitbucket.org/jaehong-blog/simple-http"
	"time"
	"strconv"
)

var port string

func init() {
	flag.StringVar(&port, "port", "0", "server port")
}

func Test_Listen(t *testing.T) {
	t.Log(port)

	s := simpleHttp.Server{}

	t.Log("server start", time.Now())

	iPort, _ := strconv.Atoi(port)
	srv := s.Listen(iPort)
	
	defer srv.Close()

	res, err := http.Get("http://localhost:" + port)
	if err != nil {
		t.Error(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Error(res)
	}
}
