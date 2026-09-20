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

// arrayItem is the element type of the JSON array request body.
type arrayItem struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// arrayItemExtra is the nested struct attribute of arrayItemNested.
type arrayItemExtra struct {
	Key1 string   `json:"key1"`
	Tags []string `json:"tags"`
}

// arrayItemNested is the element type containing nested struct attribute.
type arrayItemNested struct {
	Role    string         `json:"role"`
	Content string         `json:"content"`
	Extra   arrayItemExtra `json:"extra"`
}

// arrayReqTagged is the request struct that enables the JSON array request body by the
// `type:"array"` tag in its `g.Meta`. The JSON array request body is mapped to the first
// slice attribute `Items`, and the following slice attribute `Tags` is not affected.
type arrayReqTagged struct {
	g.Meta     `mime:"application/json" method:"post" path:"/array/tagged" type:"array"`
	Items      []arrayItem `json:"items"`
	Tags       []string    `json:"tags,omitempty"`
	ExtraField string      `json:"extraField"`
}

// arrayReqUntagged is the same request struct as arrayReqTagged but without the
// `type:"array"` tag, of which the JSON array request body is not parsed, just as before.
type arrayReqUntagged struct {
	g.Meta     `mime:"application/json" method:"post" path:"/array/untagged"`
	Items      []arrayItem `json:"items"`
	ExtraField string      `json:"extraField"`
}

// arrayReqNested is the request struct with nested struct elements.
type arrayReqNested struct {
	g.Meta `mime:"application/json" method:"post" path:"/array/nested" type:"array"`
	Items  []arrayItemNested `json:"items"`
}

// arrayRes is the response struct in the tests of the JSON array request body.
type arrayRes struct {
	ItemsCount   int    `json:"itemsCount"`
	TagsCount    int    `json:"tagsCount"`
	FirstId      int    `json:"firstId"`
	FirstContent string `json:"firstContent"`
	ExtraField   string `json:"extraField"`
}

func writeArrayRes(r *ghttp.Request, items []arrayItem, tags []string, extraField string) {
	var res = arrayRes{
		ItemsCount: len(items),
		TagsCount:  len(tags),
		ExtraField: extraField,
	}
	if len(items) > 0 {
		res.FirstId = items[0].Id
		res.FirstContent = items[0].Content
	}
	r.Response.WriteJson(res)
}

// arrayController is the strict route controller, of which the request struct fields
// are cached by the route handler.
type arrayController struct{}

func (c *arrayController) Tagged(ctx context.Context, req *arrayReqTagged) (res *arrayRes, err error) {
	writeArrayRes(g.RequestFromCtx(ctx), req.Items, req.Tags, req.ExtraField)
	return
}

func (c *arrayController) Untagged(ctx context.Context, req *arrayReqUntagged) (res *arrayRes, err error) {
	writeArrayRes(g.RequestFromCtx(ctx), req.Items, nil, req.ExtraField)
	return
}

func (c *arrayController) Nested(ctx context.Context, req *arrayReqNested) (res *arrayRes, err error) {
	var results = make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		results = append(results, fmt.Sprintf("%s/%s/%s/%v", item.Role, item.Content, item.Extra.Key1, item.Extra.Tags))
	}
	g.RequestFromCtx(ctx).Response.WriteJson(g.Map{
		"itemsCount": len(req.Items),
		"results":    results,
	})
	return
}

// rawArrayHandler is the none-strict route handler, of which the request struct fields
// are retrieved by reflection instead of the cached route handler fields.
func rawArrayHandler(r *ghttp.Request) {
	var req *arrayReqTagged
	if err := r.Parse(&req); err != nil {
		r.Response.WriteExit(err)
	}
	writeArrayRes(r, req.Items, req.Tags, req.ExtraField)
}

func Test_Params_JsonArray_RequestBody(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(arrayController))
	})
	s.BindHandler("/array/raw", rawArrayHandler)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	var (
		prefix = fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client = g.Client()
	)
	client.SetPrefix(prefix)

	var cases = []struct {
		name   string
		path   string
		body   string
		expect string
	}{
		{
			name:   "tagged request struct with multiple items",
			path:   "/array/tagged",
			body:   `[{"id":1,"name":"item1","content":"content1"},{"id":2,"name":"item2","content":"content2"}]`,
			expect: `{"itemsCount":2,"tagsCount":0,"firstId":1,"firstContent":"content1","extraField":""}`,
		},
		{
			name:   "tagged request struct with single item",
			path:   "/array/tagged",
			body:   `[{"id":100,"name":"only one","content":"content100"}]`,
			expect: `{"itemsCount":1,"tagsCount":0,"firstId":100,"firstContent":"content100","extraField":""}`,
		},
		{
			name:   "tagged request struct with empty array",
			path:   "/array/tagged",
			body:   `[]`,
			expect: `{"itemsCount":0,"tagsCount":0,"firstId":0,"firstContent":"","extraField":""}`,
		},
		{
			// The JSON object request body is still working for the tagged request struct.
			name:   "tagged request struct with object body",
			path:   "/array/tagged",
			body:   `{"items":[{"id":7,"content":"content7"}]}`,
			expect: `{"itemsCount":1,"tagsCount":0,"firstId":7,"firstContent":"content7","extraField":""}`,
		},
		{
			// It uses the request struct fields by reflection as the handler is not a strict route.
			name:   "none-strict route handler",
			path:   "/array/raw",
			body:   `[{"id":10,"content":"content10"},{"id":20,"content":"content20"}]`,
			expect: `{"itemsCount":2,"tagsCount":0,"firstId":10,"firstContent":"content10","extraField":""}`,
		},
		{
			// The JSON array request body is not parsed for the request struct without the tag.
			name:   "untagged request struct",
			path:   "/array/untagged",
			body:   `[{"id":1,"content":"content1"}]`,
			expect: `{"itemsCount":0,"tagsCount":0,"firstId":0,"firstContent":"","extraField":""}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				t.Assert(client.PostContent(ctx, c.path, c.body), c.expect)
			})
		})
	}

	// The JSON array request body works with the JSON content type either.
	t.Run("tagged request struct with json content type", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(
				client.ContentJson().PostContent(ctx, "/array/tagged", `[{"id":3,"content":"content3"}]`),
				`{"itemsCount":1,"tagsCount":0,"firstId":3,"firstContent":"content3","extraField":""}`,
			)
		})
	})
}

func Test_Params_JsonArray_RequestBodyNested(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(arrayController))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		var body = `[
			{"role":"user","content":"hello","extra":{"key1":"value1","tags":["tag1","tag2"]}},
			{"role":"assistant","content":"world","extra":{"key1":"value2","tags":["tag3"]}}
		]`
		t.Assert(
			client.PostContent(ctx, "/array/nested", body),
			`{"itemsCount":2,"results":["user/hello/value1/[tag1 tag2]","assistant/world/value2/[tag3]"]}`,
		)
	})
}
