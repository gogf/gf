// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp_test

import (
	"fmt"
	"github.com/jin502437344/gf/internal/intlog"
	"testing"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_Router_DomainBasic(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
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
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/john"), "Not Found")
		t.Assert(client.GetContent("/john/update"), "Not Found")
		t.Assert(client.GetContent("/john/edit"), "Not Found")
		t.Assert(client.GetContent("/user/list/100.html"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))
		t.Assert(client.GetContent("/john"), "")
		t.Assert(client.GetContent("/john/update"), "john")
		t.Assert(client.GetContent("/john/edit"), "edit")
		t.Assert(client.GetContent("/user/list/100.html"), "100")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))
		t.Assert(client.GetContent("/john"), "")
		t.Assert(client.GetContent("/john/update"), "john")
		t.Assert(client.GetContent("/john/edit"), "edit")
		t.Assert(client.GetContent("/user/list/100.html"), "100")
	})
}

func Test_Router_DomainMethod(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("GET:/get", func(r *ghttp.Request) {

	})
	d.BindHandler("POST:/post", func(r *ghttp.Request) {

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
		t.Assert(resp1.StatusCode, 404)

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
		t.Assert(resp4.StatusCode, 404)
	})

	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

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

	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

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

func Test_Router_DomainStatus(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
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
		t.Assert(resp1.StatusCode, 404)

		resp2, err := client.Get("/300")
		defer resp2.Close()
		t.Assert(err, nil)
		t.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/400")
		defer resp3.Close()
		t.Assert(err, nil)
		t.Assert(resp3.StatusCode, 404)

		resp4, err := client.Get("/500")
		defer resp4.Close()
		t.Assert(err, nil)
		t.Assert(resp4.StatusCode, 404)
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

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
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

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
	})
}

func Test_Router_DomainCustomStatusHandler(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	d.BindStatusHandler(404, func(r *ghttp.Request) {
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

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/ThisDoesNotExist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "hello")
		t.Assert(client.GetContent("/ThisDoesNotExist"), "404 page")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "hello")
		t.Assert(client.GetContent("/ThisDoesNotExist"), "404 page")
	})
}

func Test_Router_Domain404(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
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

		t.Assert(client.GetContent("/"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "hello")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "hello")
	})
}

func Test_Router_DomainGroup(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.Group("/", func(group *ghttp.RouterGroup) {
		group.Group("/app", func(gApp *ghttp.RouterGroup) {
			gApp.GET("/{table}/list/{page}.html", func(r *ghttp.Request) {
				intlog.Print("/{table}/list/{page}.html")
				r.Response.Write(r.Get("table"), "&", r.Get("page"))
			})
			gApp.GET("/order/info/{order_id}", func(r *ghttp.Request) {
				intlog.Print("/order/info/{order_id}")
				r.Response.Write(r.Get("order_id"))
			})
			gApp.DELETE("/comment/{id}", func(r *ghttp.Request) {
				intlog.Print("/comment/{id}")
				r.Response.Write(r.Get("id"))
			})
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client1 := ghttp.NewClient()
		client1.SetPrefix(fmt.Sprintf("http://local:%d", p))

		client2 := ghttp.NewClient()
		client2.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client1.GetContent("/app/t/list/2.html"), "t&2")
		t.Assert(client2.GetContent("/app/t/list/2.html"), "Not Found")

		t.Assert(client1.GetContent("/app/order/info/2"), "2")
		t.Assert(client2.GetContent("/app/order/info/2"), "Not Found")

		t.Assert(client1.GetContent("/app/comment/20"), "Not Found")
		t.Assert(client2.GetContent("/app/comment/20"), "Not Found")

		t.Assert(client1.DeleteContent("/app/comment/20"), "20")
		t.Assert(client2.DeleteContent("/app/comment/20"), "Not Found")
	})
}
