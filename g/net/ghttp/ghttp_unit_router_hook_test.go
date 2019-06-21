// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func Test_Router_Hook_Basic(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHookHandlerByMap("/*", map[string]ghttp.HandlerFunc{
		"BeforeServe":  func(r *ghttp.Request) { r.Response.Write("1") },
		"AfterServe":   func(r *ghttp.Request) { r.Response.Write("2") },
		"BeforeOutput": func(r *ghttp.Request) { r.Response.Write("3") },
		"AfterOutput":  func(r *ghttp.Request) { r.Response.Write("4") },
		"BeforeClose":  func(r *ghttp.Request) { r.Response.Write("5") },
		"AfterClose":   func(r *ghttp.Request) { r.Response.Write("6") },
	})
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
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
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/priority/*any", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.BindHookHandlerByMap("/priority/show", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.Write("3")
		},
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
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
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.Write("1")
		},
	})
	s.BindHookHandlerByMap("/multi-hook", map[string]ghttp.HandlerFunc{
		"BeforeServe": func(r *ghttp.Request) {
			r.Response.Write("2")
		},
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/multi-hook"), "12show")
	})
}
