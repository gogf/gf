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

func Test_Router_Group_Group(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(func(r *ghttp.Request) {
			r.Response.Write("1")
			r.Middleware.Next()
			r.Response.Write("2")
		})
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
		group.Group("/order", func(group *ghttp.RouterGroup) {
			group.GET("/list", func(r *ghttp.Request) {
				r.Response.Write("list")
			})
			group.PUT("/update", func(r *ghttp.Request) {
				r.Response.Write("update")
			})
		})
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/info", func(r *ghttp.Request) {
				r.Response.Write("info")
			})
			group.POST("/edit", func(r *ghttp.Request) {
				r.Response.Write("edit")
			})
			group.DELETE("/drop", func(r *ghttp.Request) {
				r.Response.Write("drop")
			})
		})
		group.Group("/hook", func(group *ghttp.RouterGroup) {
			group.Hook("/*", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook any")
			})
			group.Hook("/:name", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
				r.Response.Write("hook name")
			})
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/api.v2"), "Not Found")
		gtest.Assert(client.GetContent("/api.v2/test"), "1test2")
		gtest.Assert(client.GetContent("/api.v2/hook"), "hook any")
		gtest.Assert(client.GetContent("/api.v2/hook/name"), "hook namehook any")
		gtest.Assert(client.GetContent("/api.v2/hook/name/any"), "hook any")
		gtest.Assert(client.GetContent("/api.v2/order/list"), "1list2")
		gtest.Assert(client.GetContent("/api.v2/order/update"), "Not Found")
		gtest.Assert(client.PutContent("/api.v2/order/update"), "1update2")
		gtest.Assert(client.GetContent("/api.v2/user/drop"), "Not Found")
		gtest.Assert(client.DeleteContent("/api.v2/user/drop"), "1drop2")
		gtest.Assert(client.GetContent("/api.v2/user/edit"), "Not Found")
		gtest.Assert(client.PostContent("/api.v2/user/edit"), "1edit2")
		gtest.Assert(client.GetContent("/api.v2/user/info"), "1info2")
	})
}
