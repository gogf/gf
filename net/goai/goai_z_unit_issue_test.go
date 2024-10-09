// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var ctx = context.Background()

type Issue3664DefaultReq struct {
	g.Meta `path:"/default" method:"post"`
	Name   string
}
type Issue3664DefaultRes struct{}

type Issue3664RequiredTagReq struct {
	g.Meta `path:"/required-tag" required:"true" method:"post"`
	Name   string
}
type Issue3664RequiredTagRes struct{}

type Issue3664 struct{}

func (Issue3664) Default(ctx context.Context, req *Issue3664DefaultReq) (res *Issue3664DefaultRes, err error) {
	res = &Issue3664DefaultRes{}
	return
}

func (Issue3664) RequiredTag(
	ctx context.Context, req *Issue3664RequiredTagReq,
) (res *Issue3664RequiredTagRes, err error) {
	res = &Issue3664RequiredTagRes{}
	return
}

// https://github.com/gogf/gf/issues/3664
func Test_Issue3664(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3664),
			)
		})
		s.SetLogger(nil)
		s.SetOpenApiPath("/api.json")
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		apiContent := c.GetBytes(ctx, "/api.json")
		j, err := gjson.LoadJson(apiContent)
		t.AssertNil(err)
		t.Assert(j.Get(`paths./default.post.requestBody.required`).String(), "")
		t.Assert(j.Get(`paths./required-tag.post.requestBody.required`).String(), "true")
	})
}

type Issue3135DefaultReq struct {
	g.Meta `path:"/demo/colors" method:"POST" summary:"颜色 - 保存" tags:"颜色管理" description:"颜色 - 保存"`
	ID     uint64      `json:"id,string" dc:"ID" v:"id-zero"`
	Color  string      `json:"color" dc:"颜色值16进制表示法" v:"required|max-length:10"`
	Rgba   *gjson.Json `json:"rgba" dc:"颜色值RGBA表示法" v:"required|json" type:"string"`
}
type Issue3135DefaultRes struct{}

type Issue3135 struct{}

func (Issue3135) Default(ctx context.Context, req *Issue3135DefaultReq) (res *Issue3135DefaultRes, err error) {
	res = &Issue3135DefaultRes{}
	return
}

// https://github.com/gogf/gf/issues/3135
func Test_Issue3135(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3135),
			)
		})
		s.SetLogger(nil)
		s.SetOpenApiPath("/api.json")
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		var (
			api           = s.GetOpenApi()
			reqPath       = "github.com.gogf.gf.v2.net.goai_test.Issue3135DefaultReq"
			rgbType       = api.Components.Schemas.Get(reqPath).Value.Properties.Get("rgba").Value.Type
			requiredArray = api.Components.Schemas.Get(reqPath).Value.Required
		)
		t.Assert(rgbType, "string")
		t.AssertIN("rgba", requiredArray)
	})
}
