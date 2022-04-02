// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Template_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New(gtest.DataPath("template", "basic"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("index.html", g.Map{
				"name": "john",
			})
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Name:john")
		t.Assert(client.GetContent(ctx, "/"), "Name:john")
	})
}

func Test_Template_Encode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New(gtest.DataPath("template", "basic"))
		v.SetAutoEncode(true)
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("index.html", g.Map{
				"name": "john",
			})
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Name:john")
		t.Assert(client.GetContent(ctx, "/"), "Name:john")
	})
}

func Test_Template_Layout1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New(gtest.DataPath("template", "layout1"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/layout", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main1.html",
			})
			t.AssertNil(err)
		})
		s.BindHandler("/nil", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", nil)
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/layout"), "123")
		t.Assert(client.GetContent(ctx, "/nil"), "123")
	})
}

func Test_Template_Layout2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New(gtest.DataPath("template", "layout2"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/main1", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main1.html",
			})
			t.AssertNil(err)
		})
		s.BindHandler("/main2", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main2.html",
			})
			t.AssertNil(err)
		})
		s.BindHandler("/nil", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", nil)
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/main1"), "a1b")
		t.Assert(client.GetContent(ctx, "/main2"), "a2b")
		t.Assert(client.GetContent(ctx, "/nil"), "ab")
	})
}

func Test_Template_BuildInVarRequest(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.BindHandler("/:table/test", func(r *ghttp.Request) {
			err := r.Response.WriteTplContent("{{.Request.table}}")
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/user/test"), "user")
		t.Assert(client.GetContent(ctx, "/order/test"), "order")
	})
}

func Test_Template_XSS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetAutoEncode(true)
		c := "<br>"
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTplContent("{{if eq 1 1}}{{.v}}{{end}}", g.Map{
				"v": c,
			})
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), ghtml.Entities(c))
	})
}
