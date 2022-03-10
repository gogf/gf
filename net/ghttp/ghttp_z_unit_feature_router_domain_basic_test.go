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
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Router_DomainBasic(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.BindHandler("/:name", func(r *ghttp.Request) {
		r.Response.Write("/:name")
	})
	d.BindHandler("/:name/update", func(r *ghttp.Request) {
		r.Response.Write(r.Get("name"))
	})
	d.BindHandler("/:name/:action", func(r *ghttp.Request) {
		r.Response.Write(r.Get("action"))
	})
	d.BindHandler("/:name/*any", func(r *ghttp.Request) {
		r.Response.Write(r.Get("any"))
	})
	d.BindHandler("/user/list/{field}.html", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/john"), "Not Found")
		t.Assert(client.GetContent(ctx, "/john/update"), "Not Found")
		t.Assert(client.GetContent(ctx, "/john/edit"), "Not Found")
		t.Assert(client.GetContent(ctx, "/user/list/100.html"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/john"), "")
		t.Assert(client.GetContent(ctx, "/john/update"), "john")
		t.Assert(client.GetContent(ctx, "/john/edit"), "edit")
		t.Assert(client.GetContent(ctx, "/user/list/100.html"), "100")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/john"), "")
		t.Assert(client.GetContent(ctx, "/john/update"), "john")
		t.Assert(client.GetContent(ctx, "/john/edit"), "edit")
		t.Assert(client.GetContent(ctx, "/user/list/100.html"), "100")
	})
}

func Test_Router_DomainMethod(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.BindHandler("GET:/get", func(r *ghttp.Request) {

	})
	d.BindHandler("POST:/post", func(r *ghttp.Request) {

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
		t.Assert(resp1.StatusCode, 404)

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
		t.Assert(resp4.StatusCode, 404)
	})

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", s.GetListenedPort()))

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

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))

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

func Test_Router_DomainStatus(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.BindHandler("/200", func(r *ghttp.Request) {
		r.Response.WriteStatus(200)
	})
	d.BindHandler("/300", func(r *ghttp.Request) {
		r.Response.WriteStatus(300)
	})
	d.BindHandler("/400", func(r *ghttp.Request) {
		r.Response.WriteStatus(400)
	})
	d.BindHandler("/500", func(r *ghttp.Request) {
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
		t.Assert(resp1.StatusCode, 404)

		resp2, err := client.Get(ctx, "/300")
		defer resp2.Close()
		t.AssertNil(err)
		t.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get(ctx, "/400")
		defer resp3.Close()
		t.AssertNil(err)
		t.Assert(resp3.StatusCode, 404)

		resp4, err := client.Get(ctx, "/500")
		defer resp4.Close()
		t.AssertNil(err)
		t.Assert(resp4.StatusCode, 404)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", s.GetListenedPort()))

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
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))

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
	})
}

func Test_Router_DomainCustomStatusHandler(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	d.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Write("404 page")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/ThisDoesNotExist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
		t.Assert(client.GetContent(ctx, "/ThisDoesNotExist"), "404 page")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
		t.Assert(client.GetContent(ctx, "/ThisDoesNotExist"), "404 page")
	})
}

func Test_Router_Domain404(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "hello")
	})
}

func Test_Router_DomainGroup(t *testing.T) {
	s := g.Server(guid.S())
	d := s.Domain("localhost, local")
	d.Group("/", func(group *ghttp.RouterGroup) {
		group.Group("/app", func(group *ghttp.RouterGroup) {
			group.GET("/{table}/list/{page}.html", func(r *ghttp.Request) {
				intlog.Print(r.Context(), "/{table}/list/{page}.html")
				r.Response.Write(r.Get("table"), "&", r.Get("page"))
			})
			group.GET("/order/info/{order_id}", func(r *ghttp.Request) {
				intlog.Print(r.Context(), "/order/info/{order_id}")
				r.Response.Write(r.Get("order_id"))
			})
			group.DELETE("/comment/{id}", func(r *ghttp.Request) {
				intlog.Print(r.Context(), "/comment/{id}")
				r.Response.Write(r.Get("id"))
			})
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client1 := g.Client()
		client1.SetPrefix(fmt.Sprintf("http://local:%d", s.GetListenedPort()))

		client2 := g.Client()
		client2.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client1.GetContent(ctx, "/app/t/list/2.html"), "t&2")
		t.Assert(client2.GetContent(ctx, "/app/t/list/2.html"), "Not Found")

		t.Assert(client1.GetContent(ctx, "/app/order/info/2"), "2")
		t.Assert(client2.GetContent(ctx, "/app/order/info/2"), "Not Found")

		t.Assert(client1.GetContent(ctx, "/app/comment/20"), "Not Found")
		t.Assert(client2.GetContent(ctx, "/app/comment/20"), "Not Found")

		t.Assert(client1.DeleteContent(ctx, "/app/comment/20"), "20")
		t.Assert(client2.DeleteContent(ctx, "/app/comment/20"), "Not Found")
	})
}
