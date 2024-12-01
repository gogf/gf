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
	"github.com/gogf/gf/v2/errors/gcode"
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

type Issue3889CommonRes struct {
	g.Meta  `mime:"application/json"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Issue3889Req struct {
	g.Meta `path:"/default" method:"post"`
	Name   string
}
type Issue3889Res struct {
	g.Meta `status:"201" resEg:"testdata/Issue3889JsonFile/201.json"`
	Info   string `json:"info" eg:"Created!"`
}

// Example case
type Issue3889Res401 struct{}

// Override case 1
type Issue3889Res402 struct {
	g.Meta `mime:"application/json"`
}

// Override case 2
type Issue3889Res403 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Common response case
type Issue3889Res404 struct{}

var Issue3889ErrorRes = map[int][]gcode.Code{
	401: {
		gcode.New(1, "Aha, 401 - 1", nil),
		gcode.New(2, "Aha, 401 - 2", nil),
	},
}

func (r Issue3889Res) EnhanceResponseStatus() map[goai.EnhancedStatusCode]goai.EnhancedStatusType {
	Codes401 := Issue3889ErrorRes[401]
	// iterate Codes401 to generate Examples
	var Examples401 []interface{}
	for _, code := range Codes401 {
		example := Issue3889CommonRes{
			Code:    code.Code(),
			Message: code.Message(),
			Data:    nil,
		}
		Examples401 = append(Examples401, example)
	}
	return map[goai.EnhancedStatusCode]goai.EnhancedStatusType{
		401: {
			Response: Issue3889Res401{},
			Examples: Examples401,
		},
		402: {
			Response: Issue3889Res402{},
		},
		403: {
			Response: Issue3889Res403{},
		},
		404: {
			Response: Issue3889Res404{},
		},
		500: {
			Response: struct{}{},
		},
		501: {},
	}
}

type Issue3889 struct{}

func (Issue3889) Default(ctx context.Context, req *Issue3889Req) (res *Issue3889Res, err error) {
	res = &Issue3889Res{}
	return
}

// https://github.com/gogf/gf/issues/3889
func Test_Issue3889(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		openapi := s.GetOpenApi()
		openapi.Config.CommonResponse = Issue3889CommonRes{}
		openapi.Config.CommonResponseDataField = `Data`
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3889),
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
		t.AssertNE(j.Get(`paths./default.post.responses.500`).String(), "")
		t.Assert(j.Get(`paths./default.post.responses.501`).String(), "")
		// Check content
		commonResponseSchema := `{"properties":{"code":{"format":"int","type":"integer"},"data":{"properties":{},"type":"object"},"message":{"format":"string","type":"string"}},"type":"object"}`
		Status201ExamplesContent := `{"code 1":{"value":{"code":1,"data":"Good","message":"Aha, 201 - 1"}},"code 2":{"value":{"code":2,"data":"Not Bad","message":"Aha, 201 - 2"}}}`
		Status401ExamplesContent := `{"example 1":{"value":{"code":1,"data":null,"message":"Aha, 401 - 1"}},"example 2":{"value":{"code":2,"data":null,"message":"Aha, 401 - 2"}}}`
		Status402SchemaContent := `{"$ref":"#/components/schemas/github.com.gogf.gf.v2.net.goai_test.Issue3889Res402","description":""}`
		Issue3889Res403Ref := `{"$ref":"#/components/schemas/github.com.gogf.gf.v2.net.goai_test.Issue3889Res403","description":""}`

		t.Assert(j.Get(`paths./default.post.responses.201.content.application/json.examples`).String(), Status201ExamplesContent)
		t.Assert(j.Get(`paths./default.post.responses.401.content.application/json.examples`).String(), Status401ExamplesContent)
		t.Assert(j.Get(`paths./default.post.responses.402.content.application/json.schema`).String(), Status402SchemaContent)
		t.Assert(j.Get(`paths./default.post.responses.403.content.application/json.schema`).String(), Issue3889Res403Ref)
		t.Assert(j.Get(`paths./default.post.responses.404.content.application/json.schema`).String(), commonResponseSchema)
		t.Assert(j.Get(`paths./default.post.responses.500.content.application/json.schema`).String(), commonResponseSchema)

		api := s.GetOpenApi()
		reqPath := "github.com.gogf.gf.v2.net.goai_test.Issue3889Res403"
		schema := api.Components.Schemas.Get(reqPath).Value

		Issue3889Res403Schema := `{"properties":{"code":{"format":"int","type":"integer"},"message":{"format":"string","type":"string"}},"type":"object"}`
		t.Assert(schema, Issue3889Res403Schema)
	})
}

type Issue3930DefaultReq struct {
	g.Meta `path:"/user/{id}" method:"get" tags:"User" summary:"Get one user"`
	Id     int64 `v:"required" dc:"user id"`
}
type Issue3930DefaultRes struct {
	*Issue3930User `dc:"user"`
}
type Issue3930User struct {
	Id uint `json:"id"     orm:"id"     description:"user id"` // user id
}

type Issue3930 struct{}

func (Issue3930) Default(ctx context.Context, req *Issue3930DefaultReq) (res *Issue3930DefaultRes, err error) {
	res = &Issue3930DefaultRes{}
	return
}

// https://github.com/gogf/gf/issues/3930
func Test_Issue3930(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3930),
			)
		})
		s.SetLogger(nil)
		s.SetOpenApiPath("/api.json")
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		var (
			api     = s.GetOpenApi()
			reqPath = "github.com.gogf.gf.v2.net.goai_test.Issue3930DefaultRes"
		)
		t.AssertNE(api.Components.Schemas.Get(reqPath).Value.Properties.Get("id"), nil)
	})
}

type Issue3235DefaultReq struct {
	g.Meta `path:"/user/{id}" method:"get" tags:"User" summary:"Get one user"`
	Id     int64 `v:"required" dc:"user id"`
}
type Issue3235DefaultRes struct {
	Name string         `dc:"test name desc"`
	User *Issue3235User `dc:"test user desc"`
}
type Issue3235User struct {
	Id uint `json:"id"     orm:"id"     description:"user id"` // user id
}

type Issue3235 struct{}

func (Issue3235) Default(ctx context.Context, req *Issue3235DefaultReq) (res *Issue3235DefaultRes, err error) {
	res = &Issue3235DefaultRes{}
	return
}

// https://github.com/gogf/gf/issues/3235
func Test_Issue3235(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Use(ghttp.MiddlewareHandlerResponse)
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.Bind(
				new(Issue3235),
			)
		})
		s.SetLogger(nil)
		s.SetOpenApiPath("/api.json")
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		var (
			api     = s.GetOpenApi()
			reqPath = "github.com.gogf.gf.v2.net.goai_test.Issue3235DefaultRes"
		)

		t.Assert(api.Components.Schemas.Get(reqPath).Value.Properties.Get("Name").Value.Description,
			"test name desc")
		t.Assert(api.Components.Schemas.Get(reqPath).Value.Properties.Get("User").Description,
			"test user desc")
	})
}
