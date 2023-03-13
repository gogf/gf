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

type UserReq struct {
	g.Meta `path:"/user" tags:"User" method:"post" summary:"user api" title:"api title"`
	Id     int    `v:"required" d:"1"`
	Name   string `v:"required" in:"cookie"`
	Age    string `v:"required" in:"header"`
	// header,query,cookie,form
}

type UserRes struct {
	g.Meta `mime:"text/html" example:"string"`
}

var (
	User = cUser{}
)

type cUser struct{}

func (c *cUser) User(ctx context.Context, req *UserReq) (res *UserRes, err error) {
	g.RequestFromCtx(ctx).Response.WriteJson(req)
	return
}

func Test_Params_Tag(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(User)
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

func Benchmark_ParamTag(b *testing.B) {
	b.StopTimer()

	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(User)
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

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		client.PostContent(ctx, "/user", "key="+strconv.Itoa(i))
	}
}
