// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func init() {
	p := 8999
	s := g.Server(p)
	// HTTP method handlers.
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"GET: query: %d, %s",
				r.GetQuery("id").Int(),
				r.GetQuery("name").String(),
			)
		})
		group.PUT("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"PUT: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.POST("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"POST: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.DELETE("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"DELETE: form: %d, %s",
				r.GetForm("id").Int(),
				r.GetForm("name").String(),
			)
		})
		group.HEAD("/", func(r *ghttp.Request) {
			r.Response.Write("head")
		})
		group.OPTIONS("/", func(r *ghttp.Request) {
			r.Response.Write("options")
		})
	})
	// Client chaining operations handlers.
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/header", func(r *ghttp.Request) {
			r.Response.Writef(
				"Span-Id: %s, Trace-Id: %s",
				r.Header.Get("Span-Id"),
				r.Header.Get("Trace-Id"),
			)
		})
		group.ALL("/cookie", func(r *ghttp.Request) {
			r.Response.Writef(
				"SessionId: %s",
				r.Cookie.Get("SessionId"),
			)
		})
		group.ALL("/json", func(r *ghttp.Request) {
			r.Response.Writef(
				"Content-Type: %s, id: %d",
				r.Header.Get("Content-Type"),
				r.Get("id").Int(),
			)
		})
	})
	// Other testing handlers.
	s.Group("/var", func(group *ghttp.RouterGroup) {
		group.ALL("/json", func(r *ghttp.Request) {
			r.Response.Write(`{"id":1,"name":"john"}`)
		})
		group.ALL("/jsons", func(r *ghttp.Request) {
			r.Response.Write(`[{"id":1,"name":"john"}, {"id":2,"name":"smith"}]`)
		})
	})
	s.SetAccessLogEnabled(false)
	s.SetDumpRouterMap(false)
	s.SetPort(p)
	err := s.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 500)
}
