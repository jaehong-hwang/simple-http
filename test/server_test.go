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

var client *http.Client
var url string
var port string
var s shttp.Server

func init() {
	flag.StringVar(&port, "port", "8080", "server port")
	url = "http://localhost:" + port

	s = shttp.Server{}
	client = &http.Client{}

	iPort, _ := strconv.Atoi(port)
	go s.Listen(iPort)
}

func req(method string, url string, t *testing.T) {
	reqBody := strings.NewReader("<golang>really</golang>")
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	} else {
		defer res.Body.Close()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		t.Log("url:", url)
		t.Log("method:", method)
		t.Log("status:", res.StatusCode)
		t.Log("contents:", string(contents))
		t.Log("============================")
	}
}

func TestListen(t *testing.T) {
	t.Log("port:", port)
	t.Log("server start:", time.Now())
	t.Log("============================")

	req(http.MethodGet, url, t)
}

func TestRoute(t *testing.T) {
	s.Routes.AddRoute(http.MethodGet, "/test/AddRoute", func(req *http.Request) (string, error) {
		return "test", nil
	})

	req(http.MethodGet, url+"/test/AddRoute", t)
}

func TestSimpleRoute(t *testing.T) {
	path := "/test/SimpleRoute"

	s.Routes.Get(path, func(req *http.Request) (string, error) {
		return "get", nil
	})

	s.Routes.Post(path, func(req *http.Request) (string, error) {
		return "post", nil
	})

	s.Routes.Put(path, func(req *http.Request) (string, error) {
		return "put", nil
	})

	s.Routes.Delete(path, func(req *http.Request) (string, error) {
		return "Delete", nil
	})

	methods := [4]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	for _, method := range methods {
		req(method, url+path, t)
	}
}

func TestServerError(t *testing.T) {
	path := "/test/serverError"
	s.Routes.Get(path, func(req *http.Request) (string, error) {
		return "auth error!", shttp.ServerError{Status: 401}
	})

	req(http.MethodGet, url+path, t)
}
