// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/internal/tracing"
	"go.opentelemetry.io/otel"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gorilla/websocket"
)

var (
	ctx = context.TODO()
)

func Test_Client_Basic(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(url)

		t.Assert(g.Client().GetContent(ctx, ""), ``)
		t.Assert(client.GetContent(ctx, "/hello"), `hello`)

		_, err := g.Client().Post(ctx, "")
		t.AssertNE(err, nil)

		_, err = g.Client().PostForm("", nil)
		t.AssertNE(err, nil)
	})
}

func Test_Client_BasicAuth(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/auth", func(r *ghttp.Request) {
		if r.BasicAuth("admin", "admin") {
			r.Response.Write("1")
		} else {
			r.Response.Write("0")
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.BasicAuth("admin", "123").GetContent(ctx, "/auth"), "0")
		t.Assert(c.BasicAuth("admin", "admin").GetContent(ctx, "/auth"), "1")
	})
}

func Test_Client_Cookie(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/cookie", func(r *ghttp.Request) {
		r.Response.Write(r.Cookie.Get("test"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		c.SetCookie("test", "0123456789")
		t.Assert(c.PostContent(ctx, "/cookie"), "0123456789")
	})
}

func Test_Client_MapParam(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/map", func(r *ghttp.Request) {
		r.Response.Write(r.Get("test"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(c.GetContent(ctx, "/map", g.Map{"test": "1234567890"}), "1234567890")
	})
}

func Test_Client_Cookies(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/cookie", func(r *ghttp.Request) {
		r.Cookie.Set("test1", "1")
		r.Cookie.Set("test2", "2")
		r.Response.Write("ok")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp, err := c.Get(ctx, "/cookie")
		t.AssertNil(err)
		defer resp.Close()

		t.AssertNE(resp.Header.Get("Set-Cookie"), "")

		m := resp.GetCookieMap()
		t.Assert(len(m), 2)
		t.Assert(m["test1"], 1)
		t.Assert(m["test2"], 2)
		t.Assert(resp.GetCookie("test1"), 1)
		t.Assert(resp.GetCookie("test2"), 2)
	})
}

func Test_Client_Chain_Header(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/header1", func(r *ghttp.Request) {
		r.Response.Write(r.Header.Get("test1"))
	})
	s.BindHandler("/header2", func(r *ghttp.Request) {
		r.Response.Write(r.Header.Get("test2"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(c.Header(g.MapStrStr{"test1": "1234567890"}).GetContent(ctx, "/header1"), "1234567890")
		t.Assert(c.HeaderRaw("test1: 1234567890\ntest2: abcdefg").GetContent(ctx, "/header1"), "1234567890")
		t.Assert(c.HeaderRaw("test1: 1234567890\ntest2: abcdefg").GetContent(ctx, "/header2"), "abcdefg")
	})
}

func Test_Client_Chain_Context(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/context", func(r *ghttp.Request) {
		time.Sleep(1 * time.Second)
		r.Response.Write("ok")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		ctx, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)
		t.Assert(c.GetContent(ctx, "/context"), "")

		ctx, _ = context.WithTimeout(context.Background(), 2000*time.Millisecond)
		t.Assert(c.GetContent(ctx, "/context"), "ok")
	})
}

func Test_Client_Chain_Timeout(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/timeout", func(r *ghttp.Request) {
		time.Sleep(1 * time.Second)
		r.Response.Write("ok")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.Timeout(100*time.Millisecond).GetContent(ctx, "/timeout"), "")
		t.Assert(c.Timeout(2000*time.Millisecond).GetContent(ctx, "/timeout"), "ok")
	})
}

func Test_Client_Chain_ContentJson(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/json", func(r *ghttp.Request) {
		r.Response.Write(r.Get("name"), r.Get("score"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.ContentJson().PostContent(ctx, "/json", g.Map{
			"name":  "john",
			"score": 100,
		}), "john100")
		t.Assert(c.ContentJson().PostContent(ctx, "/json", `{"name":"john", "score":100}`), "john100")

		type User struct {
			Name  string `json:"name"`
			Score int    `json:"score"`
		}
		t.Assert(c.ContentJson().PostContent(ctx, "/json", User{"john", 100}), "john100")
	})
}

func Test_Client_Chain_ContentXml(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/xml", func(r *ghttp.Request) {
		r.Response.Write(r.Get("name"), r.Get("score"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.ContentXml().PostContent(ctx, "/xml", g.Map{
			"name":  "john",
			"score": 100,
		}), "john100")
		t.Assert(c.ContentXml().PostContent(ctx, "/xml", `{"name":"john", "score":100}`), "john100")

		type User struct {
			Name  string `json:"name"`
			Score int    `json:"score"`
		}
		t.Assert(c.ContentXml().PostContent(ctx, "/xml", User{"john", 100}), "john100")
	})
}

func Test_Client_Param_Containing_Special_Char(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("k1=", r.Get("k1"), "&k2=", r.Get("k2"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.PostContent(ctx, "/", "k1=MTIxMg==&k2=100"), "k1=MTIxMg==&k2=100")
		t.Assert(c.PostContent(ctx, "/", g.Map{
			"k1": "MTIxMg==",
			"k2": 100,
		}), "k1=MTIxMg==&k2=100")
	})
}

// It posts data along with file uploading.
// It does not url-encodes the parameters.
func Test_Client_File_And_Param(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/", func(r *ghttp.Request) {
		tmpPath := gfile.Temp(guid.S())
		err := gfile.Mkdir(tmpPath)
		gtest.AssertNil(err)
		defer gfile.Remove(tmpPath)

		file := r.GetUploadFile("file")
		_, err = file.Save(tmpPath)
		gtest.AssertNil(err)
		r.Response.Write(
			r.Get("json"),
			gfile.GetContents(gfile.Join(tmpPath, gfile.Basename(file.Filename))),
		)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		path := gtest.DataPath("upload", "file1.txt")
		data := g.Map{
			"file": "@file:" + path,
			"json": `{"uuid": "luijquiopm", "isRelative": false, "fileName": "test111.xls"}`,
		}
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.PostContent(ctx, "/", data), data["json"].(string)+gfile.GetContents(path))
	})
}

func Test_Client_Middleware(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	isServerHandler := false
	s.BindHandler("/", func(r *ghttp.Request) {
		isServerHandler = true
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		var (
			str1 = ""
			str2 = "resp body"
		)
		c := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		c.Use(func(c *gclient.Client, r *http.Request) (resp *gclient.Response, err error) {
			str1 += "a"
			resp, err = c.Next(r)
			if err != nil {
				return nil, err
			}
			str1 += "b"
			return
		})
		c.Use(func(c *gclient.Client, r *http.Request) (resp *gclient.Response, err error) {
			str1 += "c"
			resp, err = c.Next(r)
			if err != nil {
				return nil, err
			}
			str1 += "d"
			return
		})
		c.Use(func(c *gclient.Client, r *http.Request) (resp *gclient.Response, err error) {
			str1 += "e"
			resp, err = c.Next(r)
			if err != nil {
				return nil, err
			}
			resp.Response.Body = ioutil.NopCloser(bytes.NewBufferString(str2))
			str1 += "f"
			return
		})
		resp, err := c.Get(ctx, "/")
		t.Assert(str1, "acefdb")
		t.AssertNil(err)
		t.Assert(resp.ReadAllString(), str2)
		t.Assert(isServerHandler, true)

		// test abort, abort will not send
		var (
			str3     = ""
			abortStr = "abort request"
		)

		c = g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		c.Use(func(c *gclient.Client, r *http.Request) (resp *gclient.Response, err error) {
			str3 += "a"
			resp, err = c.Next(r)
			str3 += "b"
			return
		})
		c.Use(func(c *gclient.Client, r *http.Request) (*gclient.Response, error) {
			str3 += "c"
			return nil, gerror.New(abortStr)
		})
		c.Use(func(c *gclient.Client, r *http.Request) (resp *gclient.Response, err error) {
			str3 += "f"
			resp, err = c.Next(r)
			str3 += "g"
			return
		})
		resp, err = c.Get(ctx, "/")
		t.Assert(err, abortStr)
		t.Assert(str3, "acb")
		t.Assert(resp, nil)
	})
}

func Test_Client_Agent(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(r.UserAgent())
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client().SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		c.SetAgent("test")
		t.Assert(c.GetContent(ctx, "/"), "test")
	})
}

func Test_Client_Request_13_Dump(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.WriteJson(g.Map{"field": "test_for_response_body"})
	})
	s.BindHandler("/hello2", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.Writeln(g.Map{"field": "test_for_response_body"})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client().SetPrefix(url).ContentJson()
		r, err := client.Post(ctx, "/hello", g.Map{"field": "test_for_request_body"})
		t.AssertNil(err)
		dumpedText := r.RawRequest()
		t.Assert(gstr.Contains(dumpedText, "test_for_request_body"), true)
		dumpedText2 := r.RawResponse()
		fmt.Println(dumpedText2)
		t.Assert(gstr.Contains(dumpedText2, "test_for_response_body"), true)

		client2 := g.Client().SetPrefix(url).ContentType("text/html")
		r2, err := client2.Post(ctx, "/hello2", g.Map{"field": "test_for_request_body"})
		t.AssertNil(err)
		dumpedText3 := r2.RawRequest()
		t.Assert(gstr.Contains(dumpedText3, "test_for_request_body"), true)
		dumpedText4 := r2.RawResponse()
		t.Assert(gstr.Contains(dumpedText4, "test_for_request_body"), false)
		r2 = nil
		t.Assert(r2.RawRequest(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		response, _ := g.Client().Get(ctx, url, g.Map{
			"id":   10000,
			"name": "john",
		})
		response = nil
		t.Assert(response.RawRequest(), "")
	})

	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		response, _ := g.Client().Get(ctx, url, g.Map{
			"id":   10000,
			"name": "john",
		})
		response.RawDump()
		t.AssertGT(len(response.Raw()), 0)
	})
}

func Test_WebSocketClient(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/ws", func(r *ghttp.Request) {
		ws, err := r.WebSocket()
		if err != nil {
			r.Exit()
		}
		for {
			msgType, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if err = ws.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := gclient.NewWebSocket()
		client.Proxy = http.ProxyFromEnvironment
		client.HandshakeTimeout = time.Minute

		conn, _, err := client.Dial(fmt.Sprintf("ws://127.0.0.1:%d/ws", p), nil)
		t.AssertNil(err)
		defer conn.Close()

		msg := []byte("hello")
		err = conn.WriteMessage(websocket.TextMessage, msg)
		t.AssertNil(err)

		mt, data, err := conn.ReadMessage()
		t.AssertNil(err)
		t.Assert(mt, websocket.TextMessage)
		t.Assert(data, msg)
	})
}

func TestLoadKeyCrt(t *testing.T) {
	var (
		testCrtFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/upload/file1.txt"
		testKeyFile = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/upload/file2.txt"
		crtFile     = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.crt"
		keyFile     = gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata/server.key"
		tlsConfig   = &tls.Config{}
	)

	gtest.C(t, func(t *gtest.T) {
		tlsConfig, _ = gclient.LoadKeyCrt("crtFile", "keyFile")
		t.AssertNil(tlsConfig)

		tlsConfig, _ = gclient.LoadKeyCrt(crtFile, "keyFile")
		t.AssertNil(tlsConfig)

		tlsConfig, _ = gclient.LoadKeyCrt(testCrtFile, testKeyFile)
		t.AssertNil(tlsConfig)

		tlsConfig, _ = gclient.LoadKeyCrt(crtFile, keyFile)
		t.AssertNE(tlsConfig, nil)
	})
}

func TestClient_DoRequest(t *testing.T) {
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.WriteJson(g.Map{"field": "test_for_response_body"})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		url := fmt.Sprintf("127.0.0.1:%d/hello", p)
		resp, err := c.DoRequest(ctx, http.MethodGet, url)
		t.AssertNil(err)
		t.AssertNE(resp, nil)
		t.Assert(resp.ReadAllString(), "{\"field\":\"test_for_response_body\"}")

		resp.Response = nil
		bytes := resp.ReadAll()
		t.Assert(bytes, []byte{})
		resp.Close()
	})

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		url := "127.0.0.1:99999/hello"
		resp, err := c.DoRequest(ctx, http.MethodGet, url)
		t.AssertNil(resp.Response)
		t.AssertNE(err, nil)
	})
}

func TestClient_RequestVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			url = "http://127.0.0.1:99999/var/jsons"
		)
		varValue := g.Client().RequestVar(ctx, http.MethodGet, url)
		t.AssertNil(varValue)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string
		}
		var (
			users []User
			url   = "http://127.0.0.1:8999/var/jsons"
		)
		err := g.Client().RequestVar(ctx, http.MethodGet, url).Scan(&users)
		t.AssertNil(err)
		t.AssertNE(users, nil)
	})
}

type CustomProvider struct {
	*sdkTrace.TracerProvider
}

func NewCustomProvider() *CustomProvider {
	return &CustomProvider{
		TracerProvider: sdkTrace.NewTracerProvider(
			sdkTrace.WithIDGenerator(NewCustomIDGenerator()),
		),
	}
}

type CustomIDGenerator struct{}

func NewCustomIDGenerator() *CustomIDGenerator {
	return &CustomIDGenerator{}
}

func (id *CustomIDGenerator) NewIDs(ctx context.Context) (traceID trace.TraceID, spanID trace.SpanID) {
	return tracing.NewIDs()
}

func (id *CustomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) (spanID trace.SpanID) {
	return tracing.NewSpanID()
}

func TestClient_CustomProvider(t *testing.T) {
	//
	otel.SetTracerProvider(NewCustomProvider())

	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.WriteJson(g.Map{"field": "test_for_response_body"})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		url := fmt.Sprintf("127.0.0.1:%d/hello", p)
		resp, err := c.DoRequest(ctx, http.MethodGet, url)
		t.AssertNil(err)
		t.AssertNE(resp, nil)
		t.Assert(resp.ReadAllString(), "{\"field\":\"test_for_response_body\"}")
	})

	//otel.SetTracerProvider(provider)
}
