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

type UserTagInReq struct {
	g.Meta `path:"/user" tags:"User" method:"post" summary:"user api" title:"api title"`
	Id     int    `v:"required" d:"1"`
	Name   string `v:"required" in:"cookie"`
	Age    string `v:"required" in:"header"`
	// struct tag in:header,query,cookie,form
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
		t.Assert(client.PostContent(ctx, "/user", "name=&age=&id="), `{"Id":1,"Name":"john","Age":"18"}`)
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
