// +build linux darwin dragonfly freebsd netbsd openbsd

package greuseport_test

import (
	"fmt"
	"github.com/gogf/gf/g/net/greuseport"
    "html"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	httpServerOneResponse = "1"
	httpServerTwoResponse = "2"
)

var (
	httpServerOne = NewHTTPServer(httpServerOneResponse)
	httpServerTwo = NewHTTPServer(httpServerTwoResponse)
)

func NewHTTPServer(resp string) *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, resp)
	}))
}
func TestNewReusablePortListener(t *testing.T) {
	listenerOne, err := greuseport.Listen("tcp4", "localhost:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerOne.Close()

	listenerTwo, err := greuseport.Listen("tcp", "127.0.0.1:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerTwo.Close()

	listenerThree, err := greuseport.Listen("tcp6", "[::]:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerThree.Close()

	listenerFour, err := greuseport.Listen("tcp6", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFour.Close()

	listenerFive, err := greuseport.Listen("tcp4", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFive.Close()

	listenerSix, err := greuseport.Listen("tcp", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerSix.Close()
}

func TestListen(t *testing.T) {
	listenerOne, err := greuseport.Listen("tcp4", "localhost:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerOne.Close()

	listenerTwo, err := greuseport.Listen("tcp", "127.0.0.1:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerTwo.Close()

	listenerThree, err := greuseport.Listen("tcp6", "[::]:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerThree.Close()

	listenerFour, err := greuseport.Listen("tcp6", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFour.Close()

	listenerFive, err := greuseport.Listen("tcp4", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerFive.Close()

	listenerSix, err := greuseport.Listen("tcp", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerSix.Close()
}

func TestNewReusablePortServers(t *testing.T) {
	listenerOne, err := greuseport.Listen("tcp4", "localhost:10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerOne.Close()

	listenerTwo, err := greuseport.Listen("tcp6", ":10081")
	if err != nil {
		t.Error(err)
	}
	defer listenerTwo.Close()

	httpServerOne.Listener = listenerOne
	httpServerTwo.Listener = listenerTwo

	httpServerOne.Start()
	httpServerTwo.Start()

	// Server One — First Response
	resp1, err := http.Get(httpServerOne.URL)
	if err != nil {
		t.Error(err)
	}
	body1, err := ioutil.ReadAll(resp1.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body1) != httpServerOneResponse && string(body1) != httpServerTwoResponse {
		t.Errorf("Expected %#v or %#v, got %#v.", httpServerOneResponse, httpServerTwoResponse, string(body1))
	}

	// Server Two — First Response
	resp2, err := http.Get(httpServerTwo.URL)
	if err != nil {
		t.Error(err)
	}
	body2, err := ioutil.ReadAll(resp2.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body2) != httpServerOneResponse && string(body2) != httpServerTwoResponse {
		t.Errorf("Expected %#v or %#v, got %#v.", httpServerOneResponse, httpServerTwoResponse, string(body2))
	}

	httpServerTwo.Close()

	// Server One — Second Response
	resp3, err := http.Get(httpServerOne.URL)
	if err != nil {
		t.Error(err)
	}
	body3, err := ioutil.ReadAll(resp3.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body3) != httpServerOneResponse {
		t.Errorf("Expected %#v, got %#v.", httpServerOneResponse, string(body3))
	}

	// Server One — Third Response
	resp5, err := http.Get(httpServerOne.URL)
	if err != nil {
		t.Error(err)
	}
	body5, err := ioutil.ReadAll(resp5.Body)
	resp1.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body5) != httpServerOneResponse {
		t.Errorf("Expected %#v, got %#v.", httpServerOneResponse, string(body5))
	}

	httpServerOne.Close()
}

func BenchmarkNewReusablePortListener(b *testing.B) {
	for i := 0; i < b.N; i++ {
		listener, err := greuseport.Listen("tcp", ":10081")

		if err != nil {
			b.Error(err)
		} else {
			listener.Close()
		}
	}
}

func ExampleNewReusablePortListener() {
	listener, err := greuseport.Listen("tcp", ":8881")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	server := &http.Server{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(os.Getgid())
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))
	})

	panic(server.Serve(listener))
}
