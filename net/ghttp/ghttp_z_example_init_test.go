// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"time"
)

func init() {
	p := 8999
	s := g.Server(p)
	// HTTP method handlers.
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"GET: query: %d, %s",
				r.GetQueryInt("id"),
				r.GetQueryString("name"),
			)
		})
		group.PUT("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"PUT: form: %d, %s",
				r.GetFormInt("id"),
				r.GetFormString("name"),
			)
		})
		group.POST("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"POST: form: %d, %s",
				r.GetFormInt("id"),
				r.GetFormString("name"),
			)
		})
		group.DELETE("/", func(r *ghttp.Request) {
			r.Response.Writef(
				"DELETE: form: %d, %s",
				r.GetFormInt("id"),
				r.GetFormString("name"),
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
				r.GetInt("id"),
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
