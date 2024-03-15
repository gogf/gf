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

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type Feature3385_testHelloReq struct {
	g.Meta `GET:"/hello" `
}

type Feature3385_testDeleteReq struct {
	g.Meta `path:"/delete" `
}

type Feature3385_testAddReq struct {
	g.Meta `path:"/add" method:"put"`
}

type Feature3385_testPayReq struct {
	g.Meta `POST:"/pay" method:"put"`
}

type Feature3385_testLoginReq struct {
	g.Meta `path:"/login" method:"get,post"`
}

type Feature3385_testHelloRes struct {
	Reply string
}
type Feature3385_testAddRes struct {
	Reply string
}
type Feature3385_testDeleteRes struct {
	Reply string
}
type Feature3385_testPayRes struct {
	Reply string
}
type Feature3385_testLoginRes struct {
	Reply string
}

type testControllerFeature3385 struct{}

func (c *testControllerFeature3385) Hello(ctx context.Context, req *Feature3385_testHelloReq) (res *Feature3385_testHelloRes, err error) {
	return &Feature3385_testHelloRes{"hello"}, nil
}
func (c *testControllerFeature3385) Delete(ctx context.Context, req *Feature3385_testDeleteReq) (res *Feature3385_testDeleteRes, err error) {
	return &Feature3385_testDeleteRes{"delete"}, nil
}
func (c *testControllerFeature3385) Add(ctx context.Context, req *Feature3385_testAddReq) (res *Feature3385_testAddRes, err error) {
	return &Feature3385_testAddRes{"add"}, nil
}
func (c *testControllerFeature3385) Pay(ctx context.Context, req *Feature3385_testPayReq) (res *Feature3385_testPayRes, err error) {
	return &Feature3385_testPayRes{"pay"}, nil
}
func (c *testControllerFeature3385) Login(ctx context.Context, req *Feature3385_testLoginReq) (res *Feature3385_testLoginRes, err error) {
	return &Feature3385_testLoginRes{"login"}, nil
}

func Test_Router_Handler_Strict_WithObject_MethodUri(t *testing.T) {
	s := g.Server(guid.S())
	s.SetPort(56007)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(
			&testControllerFeature3385{},
		)
	})
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	// success
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.PutContent(ctx, "/add"), `{"code":0,"message":"","data":{"Reply":"add"}}`)
		t.Assert(client.DeleteContent(ctx, "/delete"), `{"code":0,"message":"","data":{"Reply":"delete"}}`)
		t.Assert(client.GetContent(ctx, "/hello"), `{"code":0,"message":"","data":{"Reply":"hello"}}`)
		t.Assert(client.GetContent(ctx, "/login"), `{"code":0,"message":"","data":{"Reply":"login"}}`)
		t.Assert(client.PostContent(ctx, "/login"), `{"code":0,"message":"","data":{"Reply":"login"}}`)
		t.Assert(client.PostContent(ctx, "/pay"), `{"code":0,"message":"","data":{"Reply":"pay"}}`)

		expect := `Not Found`
		add := client.GetContent(ctx, "/add")
		t.Assert(add, expect)
		t.Assert(client.DeleteContent(ctx, "/delete/1"), expect)
		t.Assert(client.DeleteContent(ctx, "/hello"), expect)
		t.Assert(client.PutContent(ctx, "/login"), expect)
		t.Assert(client.DeleteContent(ctx, "/login"), expect)
		t.Assert(client.GetContent(ctx, "/pay"), expect)

	})

}
