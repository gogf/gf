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
	"github.com/gogf/gf/v2/net/goai"
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

type Issue3747CommonRes struct {
	g.Meta  `mime:"application/json"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Issue3747Req struct {
	g.Meta `path:"/default" method:"post"`
	Name   string
}
type Issue3747Res struct {
	g.Meta `status:"201" resEg:"testdata/Issue3747JsonFile/201.json"`
	Info   string `json:"info" eg:"Created!"`
}

// Example case
type Issue3747Res401 struct {
	g.Meta `resEg:"testdata/Issue3747JsonFile/401.json"`
}

// Override case 1
type Issue3747Res402 struct {
	g.Meta `mime:"application/json"`
}

// Override case 2
type Issue3747Res403 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Common response case
type Issue3747Res404 struct{}

func (r Issue3747Res) ResponseStatusMap() map[goai.StatusCode]any {
	return map[goai.StatusCode]any{
		401: Issue3747Res401{},
		402: Issue3747Res402{},
		403: Issue3747Res403{},
		404: Issue3747Res404{},
		405: struct{}{},
		407: interface{}(nil),
		406: nil,
	}
}

type Issue3747 struct{}

func (Issue3747) Default(ctx context.Context, req *Issue3747Req) (res *Issue3747Res, err error) {
	res = &Issue3747Res{}
	return
}

// https://github.com/gogf/gf/issues/3747
func Test_Issue3747(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		openapi := s.GetOpenApi()
		openapi.Config.CommonResponse = Issue3747CommonRes{}
		openapi.Config.CommonResponseDataField = `Data`
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3747),
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
		t.Assert(j.Get(`paths./default.post.responses.200`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.201`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.401`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.402`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.403`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.404`).String(), "")
		t.AssertNE(j.Get(`paths./default.post.responses.405`).String(), "")
		t.Assert(j.Get(`paths./default.post.responses.406`).String(), "")
		t.Assert(j.Get(`paths./default.post.responses.407`).String(), "")

		// Check content
		commonResponseSchema := `{"properties":{"code":{"format":"int","type":"integer"},"data":{"properties":{},"type":"object"},"message":{"format":"string","type":"string"}},"type":"object"}`
		Status201ExamplesContent := `{"code 1":{"value":{"code":1,"data":"Good","message":"Aha, 201 - 1"}},"code 2":{"value":{"code":2,"data":"Not Bad","message":"Aha, 201 - 2"}}}`
		Status401ExamplesContent := `{"example 1":{"value":{"code":1,"data":null,"message":"Aha, 401 - 1"}},"example 2":{"value":{"code":2,"data":null,"message":"Aha, 401 - 2"}}}`
		Status402SchemaContent := `{"$ref":"#/components/schemas/github.com.gogf.gf.v2.net.goai_test.Issue3747Res402"}`
		Issue3747Res403Ref := `{"$ref":"#/components/schemas/github.com.gogf.gf.v2.net.goai_test.Issue3747Res403"}`

		t.Assert(j.Get(`paths./default.post.responses.201.content.application/json.examples`).String(), Status201ExamplesContent)
		t.Assert(j.Get(`paths./default.post.responses.401.content.application/json.examples`).String(), Status401ExamplesContent)
		t.Assert(j.Get(`paths./default.post.responses.402.content.application/json.schema`).String(), Status402SchemaContent)
		t.Assert(j.Get(`paths./default.post.responses.403.content.application/json.schema`).String(), Issue3747Res403Ref)
		t.Assert(j.Get(`paths./default.post.responses.404.content.application/json.schema`).String(), commonResponseSchema)
		t.Assert(j.Get(`paths./default.post.responses.405.content.application/json.schema`).String(), commonResponseSchema)

		api := s.GetOpenApi()
		reqPath := "github.com.gogf.gf.v2.net.goai_test.Issue3747Res403"
		schema := api.Components.Schemas.Get(reqPath).Value

		Issue3747Res403Schema := `{"properties":{"code":{"format":"int","type":"integer"},"message":{"format":"string","type":"string"}},"type":"object"}`
		t.Assert(schema, Issue3747Res403Schema)
	})
}

type Issue3889DefaultReq struct {
	g.Meta `path:"/default" method:"post"`
	Name   string
}
type Issue3889DefaultRes struct{}

type Issue3889 struct{}

func (Issue3889) Default(ctx context.Context, req *Issue3889DefaultReq) (res *Issue3889DefaultRes, err error) {
	res = &Issue3889DefaultRes{}
	return
}

func OverrideOperation(operation *goai.Operation, reqObject interface{}, resObject interface{}) {
	// Find some reqObject or resObject if needed.
	// ...
	operation.Summary = "Override summary"
	operation.Responses["201"] = operation.Responses["200"]
}

// https://github.com/gogf/gf/issues/3889
func Test_Issue3889(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3889),
			)
		})
		s.SetLogger(nil)
		s.SetOpenApiPath("/api.json")
		s.SetDumpRouterMap(false)
		goai := s.GetOpenApi()
		goai.Config.OperationOverrideHook = OverrideOperation
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		apiContent := c.GetBytes(ctx, "/api.json")
		j, err := gjson.LoadJson(apiContent)
		t.AssertNil(err)
		t.Assert(j.Get(`paths./default.post.summary`).String(), "Override summary")
		t.AssertNE(j.Get(`paths./default.post.responses.201`).String(), "")
	})
}
