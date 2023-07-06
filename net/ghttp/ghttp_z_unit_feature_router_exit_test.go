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

func Test_Router_Exit(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHookHandlerByMap("/*", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe:  func(r *ghttp.Request) { r.Response.Write("1") },
		ghttp.HookAfterServe:   func(r *ghttp.Request) { r.Response.Write("2") },
		ghttp.HookBeforeOutput: func(r *ghttp.Request) { r.Response.Write("3") },
		ghttp.HookAfterOutput:  func(r *ghttp.Request) { r.Response.Write("4") },
	})
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test-start")
		r.Exit()
		r.Response.Write("test-end")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "123")
		t.Assert(client.GetContent(ctx, "/test/test"), "1test-start23")
	})
}

func Test_Router_ExitHook(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/priority/:name", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("3")
			r.ExitHook()
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
		t.Assert(client.GetContent(ctx, "/priority/show"), "3show")
	})
}

func Test_Router_ExitAll(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/priority/:name", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[ghttp.HookName]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			r.Response.Write("3")
			r.ExitAll()
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
		t.Assert(client.GetContent(ctx, "/priority/show"), "3")
	})
}
