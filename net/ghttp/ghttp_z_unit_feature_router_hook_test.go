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

func Test_Router_Hook_Basic(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHookHandlerByMap("/*", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe:  func(r *ghttp.Request) { r.Response.Write("1") },
		ghttp.HookAfterServe:   func(r *ghttp.Request) { r.Response.Write("2") },
		ghttp.HookBeforeOutput: func(r *ghttp.Request) { r.Response.Write("3") },
		ghttp.HookAfterOutput:  func(r *ghttp.Request) { r.Response.Write("4") },
	})
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "123")
		t.Assert(client.GetContent(ctx, "/test/test"), "1test23")
	})
}

func Test_Router_Hook_Fuzzy_Router(t *testing.T) {
	s := g.Server(guid.S())
	i := 1000
	pattern1 := "/:name/info"
	s.BindHookHandlerByMap(pattern1, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.SetParam("uid", i)
			i++
		},
	})
	s.BindHandler(pattern1, func(r *ghttp.Request) {
		r.Response.Write(r.Get("uid"))
	})

	pattern2 := "/{object}/list/{page}.java"
	s.BindHookHandlerByMap(pattern2, map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeOutput: func(r *ghttp.Request) {
			r.Response.SetBuffer([]byte(
				fmt.Sprint(r.Get("object"), "&", r.Get("page"), "&", i),
			))
		},
	})
	s.BindHandler(pattern2, func(r *ghttp.Request) {
		r.Response.Write(r.Router.Uri)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/john"), "Not Found")
		t.Assert(client.GetContent(ctx, "/john/info"), "1000")
		t.Assert(client.GetContent(ctx, "/john/info"), "1001")
		t.Assert(client.GetContent(ctx, "/john/list/1.java"), "john&1&1002")
		t.Assert(client.GetContent(ctx, "/john/list/2.java"), "john&2&1002")
	})
}

func Test_Router_Hook_Priority(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/priority/:name", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("3")
		},
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/priority/show"), "312show")
		t.Assert(client.GetContent(ctx, "/priority/any/any"), "2")
		t.Assert(client.GetContent(ctx, "/priority/name"), "12")
	})
}

func Test_Router_Hook_Multi(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/multi-hook", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/multi-hook", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/multi-hook", map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Not Found")
		t.Assert(client.GetContent(ctx, "/multi-hook"), "12show")
	})
}

func Test_Router_Hook_ExitAll(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.Group("/hook", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("1")
			r.Middleware.Next()
		})
		group.ALL("/test", func(r *ghttp.Request) {
			r.Response.Write("2")
		})
	})

	s.BindHookHandler("/hook/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
		r.Response.Write("hook")
		r.ExitAll()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/test"), "test")
		t.Assert(client.GetContent(ctx, "/hook/test"), "hook")
	})
}
