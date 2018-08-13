package test

import (
	"flag"
	"io/ioutil"
	http "net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	shttp "github.com/jaehong-hwang/simple-http"
)

var port string
var s shttp.Server

func init() {
	flag.StringVar(&port, "port", "8080", "server port")

	s = shttp.Server{}

	iPort, _ := strconv.Atoi(port)
	go s.Listen(iPort)
}

func TestListen(t *testing.T) {
	t.Log("port:", port)
	t.Log("server start:", time.Now())

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
	s.AddRoute(http.MethodGet, "/test/AddRoute", func(request *http.Request) []byte {
		return []byte("test")
	})

	res, err := http.Get("http://localhost:" + port + "/test/AddRoute")
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

func TestSimpleRoute(t *testing.T) {
	path := "/test/SimpleRoute"

	s.Get(path, func(request *http.Request) []byte {
		return []byte("get")
	})

	s.Post(path, func(request *http.Request) []byte {
		return []byte("post")
	})

	s.Put(path, func(request *http.Request) []byte {
		return []byte("put")
	})

	s.Delete(path, func(request *http.Request) []byte {
		return []byte("Delete")
	})

	reqBody := strings.NewReader("<golang>really</golang>")
	client := &http.Client{}
	methods := [4]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		req, err := http.NewRequest(method, "http://localhost:"+port+path, reqBody)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil || res.StatusCode == http.StatusNotFound {
			t.Fatal(err)
		} else {
			defer res.Body.Close()
			contents, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}

			t.Log("method:", method)
			t.Log("status:", res.StatusCode)
			t.Log("contents:", string(contents))
			t.Log("============================")
		}
	}
}
