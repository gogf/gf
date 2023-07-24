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
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_OpenApi_Swagger(t *testing.T) {
	type TestReq struct {
		gmeta.Meta `method:"get" summary:"Test summary" tags:"Test"`
		Age        int
		Name       string
	}
	type TestRes struct {
		Id   int
		Age  int
		Name string
	}
	s := g.Server(guid.S())
	s.SetSwaggerPath("/swagger")
	s.SetOpenApiPath("/api.json")
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
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(c.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)

		t.Assert(gstr.Contains(c.GetContent(ctx, "/swagger/"), `API Reference`), true)
		t.Assert(gstr.Contains(c.GetContent(ctx, "/api.json"), `/test/error`), true)
	})
}

func Test_OpenApi_Multiple_Methods_Swagger(t *testing.T) {
	type TestReq struct {
		gmeta.Meta `method:"get,post" summary:"Test summary" tags:"Test"`
		Age        int
		Name       string
	}
	type TestRes struct {
		Id   int
		Age  int
		Name string
	}
	s := g.Server(guid.S())
	s.SetSwaggerPath("/swagger")
	s.SetOpenApiPath("/api.json")
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
		openapi := s.GetOpenApi()
		t.AssertNE(openapi.Paths["/test"].Get, nil)
		t.AssertNE(openapi.Paths["/test"].Post, nil)
		t.AssertNE(openapi.Paths["/test/error"].Get, nil)
		t.AssertNE(openapi.Paths["/test/error"].Post, nil)

		t.Assert(len(openapi.Paths["/test"].Get.Parameters), 2)
		t.Assert(len(openapi.Paths["/test/error"].Get.Parameters), 2)
		t.Assert(len(openapi.Components.Schemas.Get(`github.com.gogf.gf.v2.net.ghttp_test.TestReq`).Value.Properties.Map()), 2)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Only works on GET & POST methods.
		t.Assert(c.GetContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(c.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)
		t.Assert(c.PostContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(c.PostContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)

		// Not works on other methods.
		t.Assert(c.PutContent(ctx, "/test?age=18&name=john"), `{"code":65,"message":"Not Found","data":null}`)
		t.Assert(c.PutContent(ctx, "/test/error"), `{"code":65,"message":"Not Found","data":null}`)

		t.Assert(gstr.Contains(c.GetContent(ctx, "/swagger/"), `API Reference`), true)
		t.Assert(gstr.Contains(c.GetContent(ctx, "/api.json"), `/test/error`), true)
	})
}

func Test_OpenApi_Method_All_Swagger(t *testing.T) {
	type TestReq struct {
		gmeta.Meta `method:"all" summary:"Test summary" tags:"Test"`
		Age        int
		Name       string
	}
	type TestRes struct {
		Id   int
		Age  int
		Name string
	}
	s := g.Server(guid.S())
	s.SetSwaggerPath("/swagger")
	s.SetOpenApiPath("/api.json")
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
		openapi := s.GetOpenApi()
		t.AssertNE(openapi.Paths["/test"].Get, nil)
		t.AssertNE(openapi.Paths["/test"].Post, nil)
		t.AssertNE(openapi.Paths["/test"].Delete, nil)
		t.AssertNE(openapi.Paths["/test/error"].Get, nil)
		t.AssertNE(openapi.Paths["/test/error"].Post, nil)
		t.AssertNE(openapi.Paths["/test/error"].Delete, nil)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(c.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)
		t.Assert(c.PostContent(ctx, "/test?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18,"Name":"john"}}`)
		t.Assert(c.PostContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)

		t.Assert(gstr.Contains(c.GetContent(ctx, "/swagger/"), `API Reference`), true)
		t.Assert(gstr.Contains(c.GetContent(ctx, "/api.json"), `/test/error`), true)
	})
}
