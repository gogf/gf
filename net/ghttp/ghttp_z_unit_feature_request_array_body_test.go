// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// arrayBodyItem is the element type of the JSON array request body.
type arrayBodyItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// arrayBodyExtra is the nested attribute of arrayBodyNestedItem.
type arrayBodyExtra struct {
	Key  string   `json:"key"`
	Tags []string `json:"tags"`
}

// arrayBodyNestedItem is the element type containing nested structures.
type arrayBodyNestedItem struct {
	Role  string         `json:"role"`
	Extra arrayBodyExtra `json:"extra"`
}

// ArrayBodyReq declares the JSON array request body by the field tagged with `in:"body"`, of
// which the `Id` attribute is still retrieved from the request parameters as usual.
type ArrayBodyReq struct {
	g.Meta   `path:"/array-body" method:"post" mime:"application/json" summary:"array body api"`
	Id       int             `json:"id" d:"100"`
	Messages []arrayBodyItem `json:"messages" dc:"message list" in:"body"`
}

// ArrayBodyRes is the response struct showing what the handler received.
type ArrayBodyRes struct {
	Id     int      `json:"id"`
	Count  int      `json:"count"`
	IsNil  bool     `json:"isNil"`
	Names  []string `json:"names"`
	Result string   `json:"result"`
}

// ArrayBodyPathReq declares the JSON array request body together with a path parameter.
type ArrayBodyPathReq struct {
	g.Meta   `path:"/array-body/{id}" method:"post" mime:"application/json"`
	Id       int             `json:"id"`
	Messages []arrayBodyItem `json:"messages" in:"body"`
}

// ArrayBodyNestedReq declares the JSON array request body of nested elements.
type ArrayBodyNestedReq struct {
	g.Meta `path:"/array-body-nested" method:"post" mime:"application/json"`
	Items  []arrayBodyNestedItem `json:"items" in:"body"`
}

type ArrayBodyPointerReq struct {
	g.Meta `path:"/array-body-pointer" method:"post" mime:"application/json"`
	Items  *[]int `json:"items" in:"body"`
}

// ArrayBodyPlainReq is an ordinary request struct that does not declare the array request body.
type ArrayBodyPlainReq struct {
	g.Meta `path:"/array-body-plain" method:"post" mime:"application/json"`
	Name   string `json:"name"`
}

type arrayBodyPointerRes struct {
	Items *[]int `json:"items"`
}

func newArrayBodyRes(res *ArrayBodyRes, id int, items []arrayBodyItem, result string) *ArrayBodyRes {
	res.Id = id
	res.Count = len(items)
	res.IsNil = items == nil
	res.Result = result
	for _, item := range items {
		res.Names = append(res.Names, item.Name)
	}
	return res
}

// arrayBodyController is the controller for the JSON array request body tests.
type arrayBodyController struct{}

func (c *arrayBodyController) Create(ctx context.Context, req *ArrayBodyReq) (res *ArrayBodyRes, err error) {
	return newArrayBodyRes(&ArrayBodyRes{}, req.Id, req.Messages, ""), nil
}

func (c *arrayBodyController) Update(ctx context.Context, req *ArrayBodyPathReq) (res *ArrayBodyRes, err error) {
	return newArrayBodyRes(&ArrayBodyRes{}, req.Id, req.Messages, ""), nil
}

func (c *arrayBodyController) Nested(ctx context.Context, req *ArrayBodyNestedReq) (res *ArrayBodyRes, err error) {
	res = &ArrayBodyRes{
		Count: len(req.Items),
		IsNil: req.Items == nil,
	}
	for _, item := range req.Items {
		res.Names = append(res.Names, item.Role)
		res.Result = fmt.Sprintf("%s/%s/%v", item.Extra.Key, item.Role, item.Extra.Tags)
	}
	return
}

func (c *arrayBodyController) Pointer(ctx context.Context, req *ArrayBodyPointerReq) (*arrayBodyPointerRes, error) {
	return &arrayBodyPointerRes{Items: req.Items}, nil
}

func (c *arrayBodyController) Plain(ctx context.Context, req *ArrayBodyPlainReq) (res *ArrayBodyRes, err error) {
	return &ArrayBodyRes{Result: req.Name}, nil
}

// startArrayBodyServer starts a server for the tests and waits for its listening port.
func startArrayBodyServer(t *testing.T, s *ghttp.Server) *gclient.Client {
	t.Helper()
	s.SetDumpRouterMap(false)
	s.Start()
	t.Cleanup(func() {
		if err := s.Shutdown(); err != nil {
			t.Error(err)
		}
	})
	for i := 0; i < 100 && s.GetListenedPort() == 0; i++ {
		time.Sleep(10 * time.Millisecond)
	}
	if s.GetListenedPort() == 0 {
		t.Fatal("server did not start listening")
	}
	var client = g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	return client
}

func newArrayBodyServer() *ghttp.Server {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(new(arrayBodyController))
	})
	return s
}

func Test_Request_JsonArrayBody(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(
			ctx, "/array-body?some=1", `[{"id":1,"name":"john"},{"id":2,"name":"smith"}]`,
		))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.count").Int(), 2)
		t.Assert(parsed.Get("data.isNil").Bool(), false)
		t.Assert(parsed.Get("data.names").Strings(), []string{"john", "smith"})
		// The default value of the attribute is still merged as usual.
		t.Assert(parsed.Get("data.id").Int(), 100)
	})

	// The array body is accepted no matter what the content type is, just like the object body.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentType("text/plain").PostContent(
			ctx, "/array-body", `[{"id":3,"name":"anne"}]`,
		))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.count").Int(), 1)
		t.Assert(parsed.Get("data.names").Strings(), []string{"anne"})
	})
}

func Test_Request_JsonArrayBody_BodyTakesPrecedenceOverQuery(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(
			ctx, "/array-body?messages[]=99", `[{"id":1,"name":"body"}]`,
		))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.names").Strings(), []string{"body"})
	})
}

func Test_Request_JsonArrayBody_Pointer(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body-pointer", `[3,4]`))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.items").Ints(), []int{3, 4})
	})
}

func Test_Request_JsonArrayBody_WithPathParam(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(
			ctx, "/array-body/7", `[{"id":1,"name":"john"}]`,
		))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.id").Int(), 7)
		t.Assert(parsed.Get("data.count").Int(), 1)
	})
}

func Test_Request_JsonArrayBody_EmptyAndAbsent(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	// An empty JSON array results in an empty but not nil slice.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body", `[]`))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.count").Int(), 0)
		t.Assert(parsed.Get("data.isNil").Bool(), false)
	})

	// No body at all leaves the field as a nil slice.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body", ``))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.count").Int(), 0)
		t.Assert(parsed.Get("data.isNil").Bool(), true)
	})

	// The field value only comes from the request body, so the request parameters of the field
	// name and of its case/symbol variants are not accepted for it.
	for _, url := range []string{
		"/array-body?messages=99",
		"/array-body?Messages=99",
		"/array-body?MESSAGES=99",
		"/array-body?m_e_s_s_a_g_e_s=99",
	} {
		gtest.C(t, func(t *gtest.T) {
			var parsed = gjson.New(client.ContentJson().PostContent(ctx, url, ``))
			t.Assert(parsed.Get("code").Int(), 0)
			t.Assert(parsed.Get("data.count").Int(), 0)
			t.Assert(parsed.Get("data.isNil").Bool(), true)
		})
	}

	// A scalar parameter of the field name is not converted into the struct slice field either.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body?Messages=oops", ``))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.isNil").Bool(), true)
	})
}

func Test_Request_JsonArrayBody_FormContentTypeNoPollution(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	// The JSON array body is also accepted under the form content type, and it is detected
	// before the form decoding, which would otherwise cut the array into bogus parameters:
	// the `&id=777` inside the array is not decoded as the request parameter `id`.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentType("application/x-www-form-urlencoded").PostContent(
			ctx, "/array-body?id=123", `[{"id":1,"name":"https://x.com?a=1&id=777&b=2"}]`,
		))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.id").Int(), 123)
		t.Assert(parsed.Get("data.count").Int(), 1)
		t.Assert(parsed.Get("data.names").Strings(), []string{"https://x.com?a=1&id=777&b=2"})
	})

	// The form body is still reported as an invalid parameter.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentType("application/x-www-form-urlencoded").PostContent(
			ctx, "/array-body", `id=888&name=fromform`,
		))
		t.Assert(parsed.Get("code").Int(), gcode.CodeInvalidParameter.Code())
	})
}

func Test_Request_JsonArrayBody_MultipartRejected(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	const boundary = "arrayBodyBoundary"
	var newMultipartBody = func(parts ...string) string {
		var buffer = bytes.NewBuffer(nil)
		for _, part := range parts {
			buffer.WriteString("--" + boundary + "\r\n" + part + "\r\n")
		}
		buffer.WriteString("--" + boundary + "--\r\n")
		return buffer.String()
	}
	// The multipart forms are not acceptable for the array body field, no matter they carry
	// normal fields, files or nothing at all.
	for _, body := range []string{
		newMultipartBody(`Content-Disposition: form-data; name="id"` + "\r\n\r\n888"),
		newMultipartBody(
			`Content-Disposition: form-data; name="file"; filename="a.txt"` + "\r\n" +
				`Content-Type: text/plain` + "\r\n\r\nhello",
		),
		newMultipartBody(),
	} {
		gtest.C(t, func(t *gtest.T) {
			var parsed = gjson.New(client.ContentType("multipart/form-data; boundary="+boundary).
				PostContent(ctx, "/array-body", body),
			)
			t.Assert(parsed.Get("code").Int(), gcode.CodeInvalidParameter.Code())
			t.Assert(parsed.Get("data").IsNil(), true)
		})
	}
}

func Test_Request_JsonArrayBody_Nested(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body-nested", `[
			{"role":"user","extra":{"key":"key1","tags":["a","b"]}},
			{"role":"assistant","extra":{"key":"key2","tags":["c"]}}
		]`))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.count").Int(), 2)
		t.Assert(parsed.Get("data.names").Strings(), []string{"user", "assistant"})
		t.Assert(parsed.Get("data.result").String(), "key2/assistant/[c]")
	})
}

func Test_Request_JsonArrayBody_ObjectBodyRejected(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	// The object body is not acceptable for the field tagged with `in:"body"`, which is reported
	// as an invalid parameter instead of being silently bound as a one-element slice.
	for _, body := range []string{
		`{"messages":[{"id":1,"name":"john"}]}`,
		`{"id":1}`,
	} {
		gtest.C(t, func(t *gtest.T) {
			var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body", body))
			t.Assert(parsed.Get("code").Int(), gcode.CodeInvalidParameter.Code())
			t.Assert(parsed.Get("data").IsNil(), true)
		})
	}
}

func Test_Request_JsonArrayBody_PlainEndpointUnchanged(t *testing.T) {
	s := newArrayBodyServer()
	client := startArrayBodyServer(t, s)

	// An ordinary request struct is not affected: the JSON array body is still reported as an
	// invalid parameter, just the same as before.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body-plain", `[{"name":"john"}]`))
		t.Assert(parsed.Get("code").Int(), gcode.CodeInvalidParameter.Code())
		t.Assert(parsed.Get("data").IsNil(), true)
	})

	// The ordinary JSON object body still works.
	gtest.C(t, func(t *gtest.T) {
		var parsed = gjson.New(client.ContentJson().PostContent(ctx, "/array-body-plain", `{"name":"john"}`))
		t.Assert(parsed.Get("code").Int(), 0)
		t.Assert(parsed.Get("data.result").String(), "john")
	})
}
