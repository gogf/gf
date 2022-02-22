// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Router_Handler_Extended_Handler_WithObject(t *testing.T) {
	type TestReq struct {
		Age  int
		Name string
	}
	type TestRes struct {
		Id   int
		Age  int
		Name string
	}
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test", func(ctx context.Context, req *TestReq) (res *TestRes, err error) {
		return &TestRes{
			Id:   1,
			Age:  req.Age,
			Name: req.Name,
		}, nil
	})
	s.BindHandler("/test/error", func(ctx context.Context, req *TestReq) (res *TestRes, err error) {
		return &TestRes{
			Id:   1,
			Age:  req.Age,
			Name: req.Name,
		}, gerror.New("error")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(client.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":null}`)
	})
}

type TestForHandlerWithObjectAndMeta1Req struct {
	g.Meta `path:"/custom-test1" method:"get"`
	Age    int
	Name   string
}
type TestForHandlerWithObjectAndMeta1Res struct {
	Id  int
	Age int
}

type TestForHandlerWithObjectAndMeta2Req struct {
	g.Meta `path:"/custom-test2" method:"get"`
	Age    int
	Name   string
}
type TestForHandlerWithObjectAndMeta2Res struct {
	Id   int
	Name string
}

type ControllerForHandlerWithObjectAndMeta1 struct{}

func (ControllerForHandlerWithObjectAndMeta1) Test1(ctx context.Context, req *TestForHandlerWithObjectAndMeta1Req) (res *TestForHandlerWithObjectAndMeta1Res, err error) {
	return &TestForHandlerWithObjectAndMeta1Res{
		Id:  1,
		Age: req.Age,
	}, nil
}

func (ControllerForHandlerWithObjectAndMeta1) Test2(ctx context.Context, req *TestForHandlerWithObjectAndMeta2Req) (res *TestForHandlerWithObjectAndMeta2Res, err error) {
	return &TestForHandlerWithObjectAndMeta2Res{
		Id:   1,
		Name: req.Name,
	}, nil
}

type TestForHandlerWithObjectAndMeta3Req struct {
	g.Meta `path:"/custom-test3" method:"get"`
	Age    int
	Name   string
}
type TestForHandlerWithObjectAndMeta3Res struct {
	Id  int
	Age int
}

type TestForHandlerWithObjectAndMeta4Req struct {
	g.Meta `path:"/custom-test4" method:"get"`
	Age    int
	Name   string
}
type TestForHandlerWithObjectAndMeta4Res struct {
	Id   int
	Name string
}

type ControllerForHandlerWithObjectAndMeta2 struct{}

func (ControllerForHandlerWithObjectAndMeta2) Test3(ctx context.Context, req *TestForHandlerWithObjectAndMeta3Req) (res *TestForHandlerWithObjectAndMeta3Res, err error) {
	return &TestForHandlerWithObjectAndMeta3Res{
		Id:  1,
		Age: req.Age,
	}, nil
}

func (ControllerForHandlerWithObjectAndMeta2) Test4(ctx context.Context, req *TestForHandlerWithObjectAndMeta4Req) (res *TestForHandlerWithObjectAndMeta4Res, err error) {
	return &TestForHandlerWithObjectAndMeta4Res{
		Id:   1,
		Name: req.Name,
	}, nil
}
func Test_Router_Handler_Extended_Handler_WithObjectAndMeta(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", new(ControllerForHandlerWithObjectAndMeta1))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), `{"code":0,"message":"","data":null}`)
		t.Assert(client.GetContent(ctx, "/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":null}`)
	})
}

func Test_Router_Handler_Extended_Handler_Group_Bind(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(ControllerForHandlerWithObjectAndMeta1),
			new(ControllerForHandlerWithObjectAndMeta2),
		)
	})
	s.Group("/api/v2", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(ControllerForHandlerWithObjectAndMeta1),
			new(ControllerForHandlerWithObjectAndMeta2),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), `{"code":0,"message":"","data":null}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":null}`)

		t.Assert(client.GetContent(ctx, "/api/v1/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)

		t.Assert(client.GetContent(ctx, "/api/v2/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
	})
}
