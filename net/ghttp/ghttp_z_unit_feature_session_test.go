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
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Session_Cookie(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.BindHandler("/remove", func(r *ghttp.Request) {
		r.Session.Remove(r.Get("k").String())
	})
	s.BindHandler("/clear", func(r *ghttp.Request) {
		r.Session.RemoveAll()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r1, e1 := client.Get(ctx, "/set?k=key1&v=100")
		if r1 != nil {
			defer r1.Close()
		}
		t.Assert(e1, nil)
		t.Assert(r1.ReadAllString(), "")

		t.Assert(client.GetContent(ctx, "/set?k=key2&v=200"), "")

		t.Assert(client.GetContent(ctx, "/get?k=key1"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/get?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key4"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/clear"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "")
	})
}

func Test_Session_Header(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.BindHandler("/remove", func(r *ghttp.Request) {
		r.Session.Remove(r.Get("k").String())
	})
	s.BindHandler("/clear", func(r *ghttp.Request) {
		r.Session.RemoveAll()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		response, e1 := client.Get(ctx, "/set?k=key1&v=100")
		if response != nil {
			defer response.Close()
		}
		sessionId := response.GetCookie(s.GetSessionIdName())
		t.Assert(e1, nil)
		t.AssertNE(sessionId, nil)
		t.Assert(response.ReadAllString(), "")

		client.SetHeader(s.GetSessionIdName(), sessionId)

		t.Assert(client.GetContent(ctx, "/set?k=key2&v=200"), "")

		t.Assert(client.GetContent(ctx, "/get?k=key1"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/get?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key4"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/clear"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "")
	})
}

func Test_Session_StorageFile(t *testing.T) {
	sessionId := ""
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
		r.Response.Write(r.Get("k").String(), "=", r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		response, e1 := client.Get(ctx, "/set?k=key&v=100")
		if response != nil {
			defer response.Close()
		}
		sessionId = response.GetCookie(s.GetSessionIdName())
		t.Assert(e1, nil)
		t.AssertNE(sessionId, nil)
		t.Assert(response.ReadAllString(), "key=100")
	})
	time.Sleep(time.Second)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.Assert(client.GetContent(ctx, "/get?k=key"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
	})
}

func Test_Session_Custom_Id(t *testing.T) {
	var (
		sessionId = "1234567890"
		key       = "key"
		value     = "value"
		p, _      = gtcp.GetFreePort()
		s         = g.Server(p)
	)
	s.BindHandler("/id", func(r *ghttp.Request) {
		if err := r.Session.SetId(sessionId); err != nil {
			r.Response.WriteExit(err.Error())
		}
		if err := r.Session.Set(key, value); err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(r.Session.Id())
	})
	s.BindHandler("/value", func(r *ghttp.Request) {
		r.Response.WriteExit(r.Session.Get(key))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r, err := client.Get(ctx, "/id")
		t.Assert(err, nil)
		defer r.Close()
		t.Assert(r.ReadAllString(), sessionId)
		t.Assert(r.GetCookie(s.GetSessionIdName()), sessionId)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.Assert(client.GetContent(ctx, "/value"), value)
	})
}
