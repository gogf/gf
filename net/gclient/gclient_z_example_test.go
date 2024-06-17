// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	crtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
	keyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
)

func init() {
	// Default server for client.
	p := 8999
	s := g.Server(p)
	// HTTP method handlers.
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"GET: query: %d, %s",
				r.GetQuery("id").Int(),
				r.GetQuery("name").String(),
			)
		})
		group.PUT("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"PUT: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.POST("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"POST: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.DELETE("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"DELETE: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.HEAD("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"HEAD: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.PATCH("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"PATCH: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.CONNECT("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"CONNECT: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.OPTIONS("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"OPTIONS: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.TRACE("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"TRACE: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
	})
	// Client chaining operations handlers.
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/header", func(r *ghttp.Request) {
			r.Response.Writef(
				"Span-Id: %s, Trace-Id: %s",
				r.Header.Get("Span-Id"),
				r.Header.Get("Trace-Id"),
			)
		})
		group.ALL("/cookie", func(r *ghttp.Request) {
			r.Response.Writef(
				"SessionId: %s",
				r.Cookie.Get("SessionId"),
			)
		})
		group.ALL("/json", func(r *ghttp.Request) {
			r.Response.Writef(
				"Content-Type: %s, id: %d",
				r.Header.Get("Content-Type"),
				r.Get("id").Int(),
			)
		})
	})
	// Other testing handlers.
	s.Group("/var", func(group *ghttp.RouterGroup) {
		group.ALL("/json", func(r *ghttp.Request) {
			r.Response.Write(`{"id":1,"name":"john"}`)
		})
		group.ALL("/jsons", func(r *ghttp.Request) {
			r.Response.Write(`[{"id":1,"name":"john"}, {"id":2,"name":"smith"}]`)
		})
	})
	s.SetAccessLogEnabled(false)
	s.SetDumpRouterMap(false)
	s.SetPort(p)
	err := s.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 500)
}

func ExampleNew() {
	var (
		ctx    = gctx.New()
		client = gclient.New()
	)

	if r, err := client.Get(ctx, "http://127.0.0.1:8999/var/json"); err != nil {
		panic(err)
	} else {
		defer r.Close()
		fmt.Println(r.ReadAllString())
	}

	// Output:
	// {"id":1,"name":"john"}
}

func ExampleClient_Clone() {
	var (
		ctx    = gctx.New()
		client = gclient.New()
	)

	client.SetCookie("key", "value")
	cloneClient := client.Clone()

	if r, err := cloneClient.Get(ctx, "http://127.0.0.1:8999/var/json"); err != nil {
		panic(err)
	} else {
		defer r.Close()
		fmt.Println(r.ReadAllString())
	}

	// Output:
	// {"id":1,"name":"john"}
}

func fromHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func ExampleNew_MultiConn_Recommend() {
	var (
		ctx    = gctx.New()
		client = g.Client()
	)

	// controls the maximum idle(keep-alive) connections to keep per-host
	client.Transport.(*http.Transport).MaxIdleConnsPerHost = 5

	for i := 0; i < 5; i++ {
		if r, err := client.Get(ctx, "http://127.0.0.1:8999/var/json"); err != nil {
			panic(err)
		} else {
			fmt.Println(r.ReadAllString())
			r.Close()
		}
	}

	// Output:
	//{"id":1,"name":"john"}
	//{"id":1,"name":"john"}
	//{"id":1,"name":"john"}
	//{"id":1,"name":"john"}
	//{"id":1,"name":"john"}
}

func ExampleClient_Header() {
	var (
		url    = "http://127.0.0.1:8999/header"
		header = g.MapStrStr{
			"Span-Id":  "0.1",
			"Trace-Id": "123456789",
		}
	)
	content := g.Client().Header(header).PostContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	fmt.Println(content)

	// Output:
	// Span-Id: 0.1, Trace-Id: 123456789
}

func ExampleClient_HeaderRaw() {
	var (
		url       = "http://127.0.0.1:8999/header"
		headerRaw = `
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3950.0 Safari/537.36
Span-Id: 0.1
Trace-Id: 123456789
`
	)
	content := g.Client().HeaderRaw(headerRaw).PostContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	fmt.Println(content)

	// Output:
	// Span-Id: 0.1, Trace-Id: 123456789
}

func ExampleClient_Cookie() {
	var (
		url    = "http://127.0.0.1:8999/cookie"
		cookie = g.MapStrStr{
			"SessionId": "123",
		}
	)
	content := g.Client().Cookie(cookie).PostContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	fmt.Println(content)

	// Output:
	// SessionId: 123
}

func ExampleClient_ContentJson() {
	var (
		url     = "http://127.0.0.1:8999/json"
		jsonStr = `{"id":10000,"name":"john"}`
		jsonMap = g.Map{
			"id":   10000,
			"name": "john",
		}
	)
	// Post using JSON string.
	fmt.Println(g.Client().ContentJson().PostContent(ctx, url, jsonStr))
	// Post using JSON map.
	fmt.Println(g.Client().ContentJson().PostContent(ctx, url, jsonMap))

	// Output:
	// Content-Type: application/json, id: 10000
	// Content-Type: application/json, id: 10000
}

func ExampleClient_Post() {
	url := "http://127.0.0.1:8999"
	// Send with string parameter in request body.
	r1, err := g.Client().Post(ctx, url, "id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r1.Close()
	fmt.Println(r1.ReadAllString())

	// Send with map parameter.
	r2, err := g.Client().Post(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err != nil {
		panic(err)
	}
	defer r2.Close()
	fmt.Println(r2.ReadAllString())

	// Output:
	// POST: form: 10000, john
	// POST: form: 10000, john
}

func ExampleClient_PostBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().PostBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// POST: form: 10000, john
}

func ExampleClient_DeleteBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().DeleteBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// DELETE: form: 10000, john
}

func ExampleClient_HeadBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().HeadBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
}

func ExampleClient_PatchBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().PatchBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// PATCH: form: 10000, john
}

func ExampleClient_ConnectBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().ConnectBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// CONNECT: form: 10000, john
}

func ExampleClient_OptionsBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().OptionsBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// OPTIONS: form: 10000, john
}

func ExampleClient_TraceBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().TraceBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// TRACE: form: 10000, john
}

func ExampleClient_PostContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().PostContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// POST: form: 10000, john
}

func ExampleClient_PostVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().PostVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}

func ExampleClient_Get() {
	var (
		ctx = context.Background()
		url = "http://127.0.0.1:8999"
	)

	// Send with string parameter along with URL.
	r1, err := g.Client().Get(ctx, url+"?id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r1.Close()
	fmt.Println(r1.ReadAllString())

	// Send with string parameter in request body.
	r2, err := g.Client().Get(ctx, url, "id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r2.Close()
	fmt.Println(r2.ReadAllString())

	// Send with map parameter.
	r3, err := g.Client().Get(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err != nil {
		panic(err)
	}
	defer r3.Close()
	fmt.Println(r3.ReadAllString())

	// Output:
	// GET: query: 10000, john
	// GET: query: 10000, john
	// GET: query: 10000, john
}

func ExampleClient_Put() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Put(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// PUT: form: 10000, john
}

func ExampleClient_GetBytes() {
	var (
		ctx = context.Background()
		url = "http://127.0.0.1:8999"
	)
	fmt.Println(string(g.Client().GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_PutBytes() {
	var (
		ctx = context.Background()
		url = "http://127.0.0.1:8999"
	)
	fmt.Println(string(g.Client().PutBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// PUT: form: 10000, john
}

func ExampleClient_GetContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().GetContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_GetVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		user *User
		ctx  = context.Background()
		url  = "http://127.0.0.1:8999/var/json"
	)
	err := g.Client().GetVar(ctx, url).Scan(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// Output:
	// &{1 john}
}

// ExampleClient_SetProxy an example for `gclient.Client.SetProxy` method.
// please prepare two proxy server before running this example.
// http proxy server listening on `127.0.0.1:1081`
// socks5 proxy server listening on `127.0.0.1:1080`
func ExampleClient_SetProxy() {
	// connect to an http proxy server
	client := g.Client()
	client.SetProxy("http://127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	resp, err := client.Get(ctx, "http://127.0.0.1:8999")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "http://127.0.0.1:8999": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
	}
	fmt.Println(err != nil)
	resp.Close()

	// connect to an http proxy server which needs auth
	client.SetProxy("http://user:password:127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	resp, err = client.Get(ctx, "http://127.0.0.1:8999")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "http://127.0.0.1:8999": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
	}
	fmt.Println(err != nil)
	resp.Close()

	// connect to a socks5 proxy server
	client.SetProxy("socks5://127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	resp, err = client.Get(ctx, "http://127.0.0.1:8999")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "http://127.0.0.1:8999": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
	}
	fmt.Println(err != nil)
	resp.Close()

	// connect to a socks5 proxy server which needs auth
	client.SetProxy("socks5://user:password@127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	resp, err = client.Get(ctx, "http://127.0.0.1:8999")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "http://127.0.0.1:8999": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
	}
	fmt.Println(err != nil)
	resp.Close()

	// Output:
	// true
	// true
	// true
	// true
}

// ExampleClient_Proxy a chain version of example for `gclient.Client.Proxy` method.
// please prepare two proxy server before running this example.
// http proxy server listening on `127.0.0.1:1081`
// socks5 proxy server listening on `127.0.0.1:1080`
// for more details, please refer to ExampleClient_SetProxy
func ExampleClient_Proxy() {
	var (
		ctx = context.Background()
	)
	client := g.Client()
	_, err := client.Proxy("http://127.0.0.1:1081").Get(ctx, "http://127.0.0.1:8999")
	fmt.Println(err != nil)

	client2 := g.Client()
	_, err = client2.Proxy("socks5://127.0.0.1:1080").Get(ctx, "http://127.0.0.1:8999")
	fmt.Println(err != nil)

	client3 := g.Client()
	_, err = client3.Proxy("").Get(ctx, "http://127.0.0.1:8999")
	fmt.Println(err != nil)

	client4 := g.Client()
	url := "http://127.0.0.1:1081" + string([]byte{0x7f})
	_, err = client4.Proxy(url).Get(ctx, "http://127.0.0.1:8999")
	fmt.Println(err != nil)

	// Output:
	// true
	// true
	// false
	// false
}

func ExampleClient_Prefix() {
	var (
		ctx = gctx.New()
	)

	s := g.Server(guid.S())
	// HTTP method handlers.
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.GET("/v1/prefix", func(r *ghttp.Request) {
			r.Response.Write("this is v1 prefix")
		})
		group.GET("/v1/hello", func(r *ghttp.Request) {
			r.Response.Write("this is v1 hello")
		})
	})
	s.SetAccessLogEnabled(false)
	s.SetDumpRouterMap(false)
	s.Start()
	time.Sleep(time.Millisecond * 100)

	// Add Client URI Prefix
	client := g.Client().Prefix(fmt.Sprintf(
		"http://127.0.0.1:%d/api/v1/", s.GetListenedPort(),
	))

	fmt.Println(string(client.GetBytes(ctx, "prefix")))
	fmt.Println(string(client.GetBytes(ctx, "hello")))

	// Output:
	// this is v1 prefix
	// this is v1 hello
}

func ExampleClient_Retry() {
	var (
		ctx = gctx.New()
		url = "http://127.0.0.1:8999"
	)
	client := g.Client().Retry(2, time.Second)

	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_RedirectLimit() {
	var (
		ctx = gctx.New()
		url = "http://127.0.0.1:8999"
	)
	client := g.Client().RedirectLimit(1)

	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_SetBrowserMode() {
	var (
		ctx = gctx.New()
		url = "http://127.0.0.1:8999"
	)
	client := g.Client().SetBrowserMode(true)

	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_SetHeader() {
	var (
		ctx = gctx.New()
		url = "http://127.0.0.1:8999"
	)
	client := g.Client()
	client.SetHeader("Server", "GoFrameServer")
	client.SetHeader("Client", "g.Client()")

	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_SetRedirectLimit() {
	go func() {
		s := g.Server()
		s.BindHandler("/hello", func(r *ghttp.Request) {
			r.Response.Writeln("hello world")
		})
		s.BindHandler("/back", func(r *ghttp.Request) {
			r.Response.RedirectBack()
		})
		s.SetDumpRouterMap(false)
		s.SetPort(8199)
		s.Run()
	}()

	time.Sleep(time.Second)

	var (
		ctx      = gctx.New()
		urlHello = "http://127.0.0.1:8199/hello"
		urlBack  = "http://127.0.0.1:8199/back"
	)
	client := g.Client().SetRedirectLimit(1)
	client.SetHeader("Referer", urlHello)

	resp, err := client.DoRequest(ctx, http.MethodGet, urlBack, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err == nil {
		fmt.Println(resp.ReadAllString())
		resp.Close()
	}

	client.SetRedirectLimit(2)
	resp, err = client.DoRequest(ctx, http.MethodGet, urlBack, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err == nil {
		fmt.Println(resp.ReadAllString())
		resp.Close()
	}

	// Output:
	// Found
	// hello world
}

func ExampleClient_SetTLSKeyCrt() {
	var (
		ctx         = gctx.New()
		url         = "http://127.0.0.1:8999"
		testCrtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/upload/file1.txt"
		testKeyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/upload/file2.txt"
	)
	client := g.Client()
	client.SetTLSKeyCrt(testCrtFile, testKeyFile)
	client.SetTLSKeyCrt(crtFile, keyFile)
	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_SetTLSConfig() {
	var (
		ctx       = gctx.New()
		url       = "http://127.0.0.1:8999"
		tlsConfig = &tls.Config{}
	)
	client := g.Client()
	client.SetTLSConfig(tlsConfig)
	fmt.Println(string(client.GetBytes(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_PutContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().PutContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// PUT: form: 10000, john
}

func ExampleClient_DeleteContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().DeleteContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// DELETE: form: 10000, john
}

func ExampleClient_HeadContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().HeadContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
}

func ExampleClient_PatchContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().PatchContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// PATCH: form: 10000, john
}

func ExampleClient_ConnectContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().ConnectContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// CONNECT: form: 10000, john
}

func ExampleClient_OptionsContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().OptionsContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// OPTIONS: form: 10000, john
}

func ExampleClient_TraceContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().TraceContent(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// TRACE: form: 10000, john
}

func ExampleClient_RequestContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().RequestContent(ctx, http.MethodGet, url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_RawRequest() {
	url := "http://127.0.0.1:8999"
	response, _ := g.Client().Get(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	fmt.Println(len(response.RawResponse()) > 100)

	// Output:
	// true
}

func ExampleClient_Delete() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Delete(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// DELETE: form: 10000, john
}

func ExampleClient_Head() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Head(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	//
}

func ExampleClient_Patch() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Patch(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// PATCH: form: 10000, john
}

func ExampleClient_Connect() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Connect(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// CONNECT: form: 10000, john
}

func ExampleClient_Options() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Options(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// OPTIONS: form: 10000, john
}

func ExampleClient_Trace() {
	url := "http://127.0.0.1:8999"
	r, _ := g.Client().Trace(ctx, url, g.Map{
		"id":   10000,
		"name": "john",
	})
	defer r.Close()
	fmt.Println(r.ReadAllString())

	// Output:
	// TRACE: form: 10000, john
}

func ExampleClient_PutVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		user *User
		ctx  = context.Background()
		url  = "http://127.0.0.1:8999/var/json"
	)
	err := g.Client().PutVar(ctx, url).Scan(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// Output:
	// &{1 john}
}

func ExampleClient_DeleteVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().DeleteVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}

func ExampleClient_HeadVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().HeadVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// []
}

func ExampleClient_PatchVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().PatchVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}

func ExampleClient_ConnectVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().ConnectVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}

func ExampleClient_OptionsVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().OptionsVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}

func ExampleClient_TraceVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().TraceVar(ctx, url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}
