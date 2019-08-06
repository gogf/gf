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

func Test_Router_Hook_Basic(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHookHandlerByMap("/*", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE:  func(r *ghttp.Request) { r.Response.Write("1") },
		ghttp.HOOK_AFTER_SERVE:   func(r *ghttp.Request) { r.Response.Write("2") },
		ghttp.HOOK_BEFORE_OUTPUT: func(r *ghttp.Request) { r.Response.Write("3") },
		ghttp.HOOK_AFTER_OUTPUT:  func(r *ghttp.Request) { r.Response.Write("4") },
	})
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(200 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "123")
		gtest.Assert(client.GetContent("/test/test"), "1test23")
	})
}

func Test_Router_Hook_Priority(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/priority/show", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/priority/:name", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.Response.Write("3")
		},
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(200 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/priority/show"), "312show")
		gtest.Assert(client.GetContent("/priority/any/any"), "2")
		gtest.Assert(client.GetContent("/priority/name"), "12")
	})
}

func Test_Router_Hook_Multi(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/multi-hook", func(r *ghttp.Request) {
		r.Response.Write("show")
	})

	s.BindHookHandlerByMap("/multi-hook", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/multi-hook", map[string]ghttp.HandlerFunc{
		ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(200 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/multi-hook"), "12show")
	})
}
