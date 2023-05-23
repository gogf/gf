// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Router_Basic1(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/:name", func(r *ghttp.Request) {
		r.Response.Write("/:name")
	})
	s.BindHandler("/:name/update", func(r *ghttp.Request) {
		r.Response.Write(r.Get("name"))
	})
	s.BindHandler("/:name/:action", func(r *ghttp.Request) {
		r.Response.Write(r.Get("action"))
	})
	s.BindHandler("/:name/*any", func(r *ghttp.Request) {
		r.Response.Write(r.Get("any"))
	})
	s.BindHandler("/user/list/{field}.html", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/john"), "")
		t.Assert(client.GetContent(ctx, "/john/update"), "john")
		t.Assert(client.GetContent(ctx, "/john/edit"), "edit")
		t.Assert(client.GetContent(ctx, "/user/list/100.html"), "100")
	})
}

func Test_Router_Basic2(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/{hash}", func(r *ghttp.Request) {
		r.Response.Write(r.Get("hash"))
	})
	s.BindHandler("/{hash}.{type}", func(r *ghttp.Request) {
		r.Response.Write(r.Get("type"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/data"), "data")
		t.Assert(client.GetContent(ctx, "/data.json"), "json")
	})
}

func Test_Router_Value(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouterMap()["hash"])
	})
	s.BindHandler("/GetRouter", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouter("name", "john").String())
	})
	s.BindHandler("/{hash}", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouter("hash").String())
	})
	s.BindHandler("/{hash}.{type}", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouter("type").String())
	})
	s.BindHandler("/{hash}.{type}.map", func(r *ghttp.Request) {
		r.Response.Write(r.GetRouterMap()["type"])
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/"), "")
		t.Assert(client.GetContent(ctx, "/GetRouter"), "john")
		t.Assert(client.GetContent(ctx, "/data"), "data")
		t.Assert(client.GetContent(ctx, "/data.json"), "json")
		t.Assert(client.GetContent(ctx, "/data.json.map"), "json")
	})
}

// HTTP method register.
func Test_Router_Method(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("GET:/get", func(r *ghttp.Request) {

	})
	s.BindHandler("POST:/post", func(r *ghttp.Request) {

	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		resp1, err := client.Get(ctx, "/get")
		defer resp1.Close()
		t.AssertNil(err)
		t.Assert(resp1.StatusCode, 200)

		resp2, err := client.Post(ctx, "/get")
		defer resp2.Close()
		t.AssertNil(err)
		t.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get(ctx, "/post")
		defer resp3.Close()
		t.AssertNil(err)
		t.Assert(resp3.StatusCode, 404)

		resp4, err := client.Post(ctx, "/post")
		defer resp4.Close()
		t.AssertNil(err)
		t.Assert(resp4.StatusCode, 200)
	})
}

// Extra char '/' of the router.
func Test_Router_ExtraChar(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/api/test"), "test")
		t.Assert(client.GetContent(ctx, "/api/test/"), "test")
		t.Assert(client.GetContent(ctx, "/api/test//"), "test")
		t.Assert(client.GetContent(ctx, "//api/test//"), "test")
		t.Assert(client.GetContent(ctx, "//api//test//"), "test")
		t.Assert(client.GetContent(ctx, "///api///test///"), "test")
	})
}

// Custom status handler.
func Test_Router_Status(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/200", func(r *ghttp.Request) {
		r.Response.WriteStatus(200)
	})
	s.BindHandler("/300", func(r *ghttp.Request) {
		r.Response.WriteStatus(300)
	})
	s.BindHandler("/400", func(r *ghttp.Request) {
		r.Response.WriteStatus(400)
	})
	s.BindHandler("/500", func(r *ghttp.Request) {
		r.Response.WriteStatus(500)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		resp1, err := client.Get(ctx, "/200")
		defer resp1.Close()
		t.AssertNil(err)
		t.Assert(resp1.StatusCode, 200)

		resp2, err := client.Get(ctx, "/300")
		defer resp2.Close()
		t.AssertNil(err)
		t.Assert(resp2.StatusCode, 300)

		resp3, err := client.Get(ctx, "/400")
		defer resp3.Close()
		t.AssertNil(err)
		t.Assert(resp3.StatusCode, 400)

		resp4, err := client.Get(ctx, "/500")
		defer resp4.Close()
		t.AssertNil(err)
		t.Assert(resp4.StatusCode, 500)

		resp5, err := client.Get(ctx, "/404")
		defer resp5.Close()
		t.AssertNil(err)
		t.Assert(resp5.StatusCode, 404)
	})
}

func Test_Router_CustomStatusHandler(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Write("404 page")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
		resp, err := client.Get(ctx, "/ThisDoesNotExist")
		defer resp.Close()
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 404)
		t.Assert(resp.ReadAllString(), "404 page")
	})
}

// 404 not found router.
func Test_Router_404(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
		resp, err := client.Get(ctx, "/ThisDoesNotExist")
		defer resp.Close()
		t.AssertNil(err)
		t.Assert(resp.StatusCode, 404)
	})
}

func Test_Router_Priority(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/admin", func(r *ghttp.Request) {
		r.Response.Write("admin")
	})
	s.BindHandler("/admin-{page}", func(r *ghttp.Request) {
		r.Response.Write("admin-{page}")
	})
	s.BindHandler("/admin-goods", func(r *ghttp.Request) {
		r.Response.Write("admin-goods")
	})
	s.BindHandler("/admin-goods-{page}", func(r *ghttp.Request) {
		r.Response.Write("admin-goods-{page}")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/admin"), "admin")
		t.Assert(client.GetContent(ctx, "/admin-1"), "admin-{page}")
		t.Assert(client.GetContent(ctx, "/admin-goods"), "admin-goods")
		t.Assert(client.GetContent(ctx, "/admin-goods-2"), "admin-goods-{page}")
	})
}
