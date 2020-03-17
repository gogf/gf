// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/os/gview"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Template_Layout1(t *testing.T) {
	gtest.Case(t, func() {
		v := gview.New(gfile.Join(gdebug.TestDataPath(), "template", "layout1"))
		p := ports.PopRand()
		s := g.Server(p)
		s.SetView(v)
		s.BindHandler("/layout", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main1.html",
			})
			gtest.Assert(err, nil)
		})
		s.BindHandler("/nil", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", nil)
			gtest.Assert(err, nil)
		})
		s.SetDumpRouterMap(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/layout"), "123")
		gtest.Assert(client.GetContent("/nil"), "123")
	})
}

func Test_Template_Layout2(t *testing.T) {
	gtest.Case(t, func() {
		v := gview.New(gfile.Join(gdebug.TestDataPath(), "template", "layout2"))
		p := ports.PopRand()
		s := g.Server(p)
		s.SetView(v)
		s.BindHandler("/main1", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main1.html",
			})
			gtest.Assert(err, nil)
		})
		s.BindHandler("/main2", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", g.Map{
				"mainTpl": "main/main2.html",
			})
			gtest.Assert(err, nil)
		})
		s.BindHandler("/nil", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("layout.html", nil)
			gtest.Assert(err, nil)
		})
		s.SetDumpRouterMap(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/main1"), "a1b")
		gtest.Assert(client.GetContent("/main2"), "a2b")
		gtest.Assert(client.GetContent("/nil"), "ab")
	})
}

func Test_Template_XSS(t *testing.T) {
	gtest.Case(t, func() {
		v := gview.New()
		v.SetAutoEncode(true)
		c := "<br>"
		p := ports.PopRand()
		s := g.Server(p)
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTplContent("{{if eq 1 1}}{{.v}}{{end}}", g.Map{
				"v": c,
			})
			gtest.Assert(err, nil)
		})
		s.SetDumpRouterMap(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), ghtml.Entities(c))
	})
}
