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

func Test_Client_Basic(t *testing.T) {
	p := ports.PopRand()
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
		client := ghttp.NewClient()
		client.SetPrefix(url)

		t.Assert(ghttp.GetContent(""), ``)
		t.Assert(client.GetContent("/hello"), `hello`)

		_, err := ghttp.Post("")
		t.AssertNE(err, nil)
	})
}

func Test_Client_Cookie(t *testing.T) {
	p := ports.PopRand()
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
		c := ghttp.NewClient()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		c.SetCookie("test", "0123456789")
		t.Assert(c.PostContent("/cookie"), "0123456789")
	})
}

func Test_Client_Cookies(t *testing.T) {
	p := ports.PopRand()
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
		c := ghttp.NewClient()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp, err := c.Get("/cookie")
		t.Assert(err, nil)
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
