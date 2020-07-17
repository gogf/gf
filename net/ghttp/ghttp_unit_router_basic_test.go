// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_Router_Basic1(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
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
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/john"), "")
		t.Assert(client.GetContent("/john/update"), "john")
		t.Assert(client.GetContent("/john/edit"), "edit")
		t.Assert(client.GetContent("/user/list/100.html"), "100")
	})
}

func Test_Router_Basic2(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/{hash}", func(r *ghttp.Request) {
		r.Response.Write(r.Get("hash"))
	})
	s.BindHandler("/{hash}.{type}", func(r *ghttp.Request) {
		r.Response.Write(r.Get("type"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/data"), "data")
		t.Assert(client.GetContent("/data.json"), "json")
	})
}

// HTTP method register.
func Test_Router_Method(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("GET:/get", func(r *ghttp.Request) {

	})
	s.BindHandler("POST:/post", func(r *ghttp.Request) {

	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp1, err := client.Get("/get")
		defer resp1.Close()
		t.Assert(err, nil)
		t.Assert(resp1.StatusCode, 200)

		resp2, err := client.Post("/get")
		defer resp2.Close()
		t.Assert(err, nil)
		t.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/post")
		defer resp3.Close()
		t.Assert(err, nil)
		t.Assert(resp3.StatusCode, 404)

		resp4, err := client.Post("/post")
		defer resp4.Close()
		t.Assert(err, nil)
		t.Assert(resp4.StatusCode, 200)
	})
}

// Extra char '/' of the router.
func Test_Router_ExtraChar(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/api/test"), "test")
		t.Assert(client.GetContent("/api/test/"), "test")
		t.Assert(client.GetContent("/api/test//"), "test")
	})
}

// Custom status handler.
func Test_Router_Status(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
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
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp1, err := client.Get("/200")
		defer resp1.Close()
		t.Assert(err, nil)
		t.Assert(resp1.StatusCode, 200)

		resp2, err := client.Get("/300")
		defer resp2.Close()
		t.Assert(err, nil)
		t.Assert(resp2.StatusCode, 300)

		resp3, err := client.Get("/400")
		defer resp3.Close()
		t.Assert(err, nil)
		t.Assert(resp3.StatusCode, 400)

		resp4, err := client.Get("/500")
		defer resp4.Close()
		t.Assert(err, nil)
		t.Assert(resp4.StatusCode, 500)

		resp5, err := client.Get("/404")
		defer resp5.Close()
		t.Assert(err, nil)
		t.Assert(resp5.StatusCode, 404)
	})
}

func Test_Router_CustomStatusHandler(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Write("404 page")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "hello")
		resp, err := client.Get("/ThisDoesNotExist")
		defer resp.Close()
		t.Assert(err, nil)
		t.Assert(resp.StatusCode, 404)
		t.Assert(resp.ReadAllString(), "404 page")
	})
}

// 404 not found router.
func Test_Router_404(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "hello")
		resp, err := client.Get("/ThisDoesNotExist")
		defer resp.Close()
		t.Assert(err, nil)
		t.Assert(resp.StatusCode, 404)
	})
}

func Test_Router_Priority(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
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
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/admin"), "admin")
		t.Assert(client.GetContent("/admin-1"), "admin-{page}")
		t.Assert(client.GetContent("/admin-goods"), "admin-goods")
		t.Assert(client.GetContent("/admin-goods-2"), "admin-goods-{page}")
	})
}
