// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// UserTagInReq struct tag "in" supports: header, cookie
type UserTagInReq struct {
	g.Meta `path:"/user" tags:"User" method:"post" summary:"user api" title:"api title"`
	Id     int    `v:"required" d:"1"`
	Name   string `v:"required" in:"cookie"`
	Age    string `v:"required" in:"header"`
}

type UserTagInRes struct {
	g.Meta `mime:"text/html" example:"string"`
}

var (
	UserTagIn = cUserTagIn{}
)

type cUserTagIn struct{}

func (c *cUserTagIn) User(ctx context.Context, req *UserTagInReq) (res *UserTagInRes, err error) {
	g.RequestFromCtx(ctx).Response.WriteJson(req)
	return
}

func Test_ParamsTagIn(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(UserTagIn)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		client.SetCookie("name", "john")
		client.SetHeader("age", "18")

		t.Assert(client.PostContent(ctx, "/user"), `{"Id":1,"Name":"john","Age":"18"}`)
		t.Assert(client.PostContent(ctx, "/user", "name=&age="), `{"Id":1,"Name":"john","Age":"18"}`)
	})
}

type UserTagDefaultReq struct {
	g.Meta   `path:"/user-default" method:"post,get" summary:"user default tag api"`
	Id       int     `v:"required" d:"1"`
	Name     string  `d:"john"`
	Age      int     `d:"18"`
	Score    float64 `d:"99.9"`
	IsVip    bool    `d:"true"`
	NickName string  `p:"nickname" d:"nickname-default"`
	EmptyStr string  `d:""`
	Email    string
	Address  string
}

type UserTagDefaultRes struct {
	g.Meta `mime:"application/json" example:"string"`
}

var (
	UserTagDefault = cUserTagDefault{}
)

type cUserTagDefault struct{}

func (c *cUserTagDefault) User(ctx context.Context, req *UserTagDefaultReq) (res *UserTagDefaultRes, err error) {
	g.RequestFromCtx(ctx).Response.WriteJson(req)
	return
}

func Test_ParamsTagDefault(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(UserTagDefault)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		// Test with no parameters, should use all default values
		resp := client.GetContent(ctx, "/user-default")
		t.Assert(resp, `{"Id":1,"Name":"john","Age":18,"Score":99.9,"IsVip":true,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)

		// Test with partial parameters (query method), should use partial default values
		resp = client.GetContent(ctx, "/user-default?id=100&name=smith")
		t.Assert(resp, `{"Id":100,"Name":"smith","Age":18,"Score":99.9,"IsVip":true,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)

		// Test with partial parameters (query method), should use partial default values
		resp = client.GetContent(ctx, "/user-default?id=100&name=smith&age")
		t.Assert(resp, `{"Id":100,"Name":"smith","Age":18,"Score":99.9,"IsVip":true,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)

		// Test providing partial parameters via POST form
		resp = client.PostContent(ctx, "/user-default", "id=200&age=30&nickname=jack")
		t.Assert(resp, `{"Id":200,"Name":"john","Age":30,"Score":99.9,"IsVip":true,"NickName":"jack","EmptyStr":"","Email":"","Address":""}`)

		// Test providing partial parameters via POST JSON
		resp = client.ContentJson().PostContent(ctx, "/user-default", g.Map{
			"id":      300,
			"name":    "bob",
			"score":   88.8,
			"address": "beijing",
		})
		t.Assert(resp, `{"Id":300,"Name":"bob","Age":18,"Score":88.8,"IsVip":true,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":"beijing"}`)

		// Test providing JSON content via GET request
		resp = client.ContentJson().PostContent(ctx, "/user-default", `{"id":500,"isVip":false}`)
		t.Assert(resp, `{"Id":500,"Name":"john","Age":18,"Score":99.9,"IsVip":false,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)

		// Test providing empty values, should use default values
		resp = client.PostContent(ctx, "/user-default", "id=400&name=&age=")
		t.Assert(resp, `{"Id":400,"Name":"","Age":0,"Score":99.9,"IsVip":true,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)

		// Test providing JSON content via GET request
		resp = client.ContentJson().GetContent(ctx, "/user-default", `{"id":500,"isVip":false}`)
		t.Assert(resp, `{"Id":500,"Name":"john","Age":18,"Score":99.9,"IsVip":false,"NickName":"nickname-default","EmptyStr":"","Email":"","Address":""}`)
	})
}

func Benchmark_ParamTagIn(b *testing.B) {
	b.StopTimer()

	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(UserTagIn)
	})
	s.SetDumpRouterMap(false)
	s.SetAccessLogEnabled(false)
	s.SetErrorLogEnabled(false)
	s.Start()
	defer s.Shutdown()
	prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	client := g.Client()
	client.SetPrefix(prefix)
	client.SetCookie("name", "john")
	client.SetHeader("age", "18")

	b.StartTimer()
	for i := 1; i < b.N; i++ {
		client.PostContent(ctx, "/user", "id="+strconv.Itoa(i))
	}
}
