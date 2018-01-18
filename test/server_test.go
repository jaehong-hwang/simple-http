package test

import (
	"testing"
	"flag"
	"net/http"
	simpleHttp "bitbucket.org/jaehong-blog/simple-http"
	"time"
	"strconv"
	"io/ioutil"
)

var port string
var s simpleHttp.Server

func init() {
	flag.StringVar(&port, "port", "0", "server port")
}

func TestListen(t *testing.T) {
	t.Log(port)

	s = simpleHttp.Server{}

	t.Log("server start", time.Now())

	iPort, _ := strconv.Atoi(port)
	go s.Listen(iPort)

	res, err := http.Get("http://localhost:" + port)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("response code: ", res.StatusCode)
	t.Log("response body: ", string(data))
}

func TestRoute(t *testing.T) {
	s.AddRoute(http.MethodGet, "/test", func(request *http.Request) []byte {
		return []byte("test")
	})

	res, err := http.Get("http://localhost:" + port + "/test")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("response code: ", res.StatusCode)
	t.Log("response body: ", string(data))
}
