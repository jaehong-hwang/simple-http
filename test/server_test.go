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
		t.Fatal(err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("response body: ", string(data))

	if res.StatusCode != http.StatusOK {
		t.Error(res)
	}
}
