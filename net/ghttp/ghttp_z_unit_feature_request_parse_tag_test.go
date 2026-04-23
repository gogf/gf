// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type parseTagCustomBindReq struct {
	g.Meta `path:"/parse-tag-custom" method:"post"`
	Name   string `json:"name" parse:"trim-space|wrap-custom:ok" v:"required"`
}

type parseTagCustomBindRes struct {
	Name string `json:"name"`
}

type parseTagCustomController struct{}

func (c *parseTagCustomController) ParseTagCustom(
	ctx context.Context, req *parseTagCustomBindReq,
) (res *parseTagCustomBindRes, err error) {
	return &parseTagCustomBindRes{Name: req.Name}, nil
}

func Test_Params_ParseTag_BuiltInAndValidation(t *testing.T) {
	type Profile struct {
		Nick string `json:"nick" parse:"trim-space|upper"`
	}
	type Req struct {
		Title     string   `json:"title" parse:"trim-space|lower" v:"required"`
		Slug      string   `json:"slug" parse:"trim-prefix:pre-|trim-suffix:-suf|upper"`
		Trimmed   string   `json:"trimmed" parse:"trim:*"`
		LeftRight string   `json:"left_right" parse:"trim-left:-|trim-right:_"`
		Sentence  string   `json:"sentence" parse:"squash-space|title"`
		Compact   string   `json:"compact" parse:"remove-space"`
		Replaced  string   `json:"replaced" parse:"replace:foo,bar"`
		Alias     *string  `json:"alias" parse:"trim-space|empty-to-nil"`
		Tags      []string `json:"tags" parse:"foreach|trim-space|upper"`
		Profile   Profile  `json:"profile"`
	}

	s := g.Server(guid.S())
	s.BindHandler("/parse-tag-built-in", func(r *ghttp.Request) {
		var req *Req
		if err := r.Parse(&req); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteJsonExit(req)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(
			client.ContentJson().PostContent(ctx, "/parse-tag-built-in", g.Map{
				"title":      "  Demo  ",
				"slug":       "pre-demo-suf",
				"trimmed":    "***demo***",
				"left_right": "---demo___",
				"sentence":   "  hello   world  ",
				"compact":    "  a \t b \n c  ",
				"replaced":   "foo foo",
				"alias":      "   ",
				"tags":       []string{"  alpha  ", " beta "},
				"profile":    g.Map{"nick": "  john  "},
			}),
			`{"title":"demo","slug":"DEMO","trimmed":"demo","left_right":"demo","sentence":"Hello World","compact":"abc","replaced":"bar bar","alias":null,"tags":["ALPHA","BETA"],"profile":{"nick":"JOHN"}}`,
		)
		t.Assert(
			client.ContentJson().PostContent(ctx, "/parse-tag-built-in", g.Map{
				"title": "   ",
			}),
			`The title field is required`,
		)
	})
}

func Test_Params_ParseTag_CustomRuleAndServiceBinding(t *testing.T) {
	ghttp.RegisterParseRule("wrap-custom", func(ctx context.Context, in ghttp.ParseFuncInput) (any, error) {
		value, ok := in.Value.(string)
		if !ok {
			return in.Value, nil
		}
		return fmt.Sprintf("[%s:%s]", strings.TrimSpace(value), in.Pattern), nil
	})
	defer ghttp.DeleteParseRule("wrap-custom")

	var controller parseTagCustomController

	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(&controller)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(
			client.ContentJson().PostContent(ctx, "/parse-tag-custom", g.Map{
				"name": "  john  ",
			}),
			`{"code":0,"message":"OK","data":{"name":"[john:ok]"}}`,
		)
	})
}

func Test_Params_ParseTag_TopLevelStructSlice(t *testing.T) {
	type Item struct {
		Title string   `json:"title" parse:"trim-space" v:"required"`
		Tags  []string `json:"tags" parse:"foreach|trim-space|lower"`
	}

	s := g.Server(guid.S())
	s.BindHandler("/parse-tag-array", func(r *ghttp.Request) {
		var items []*Item
		if err := r.Parse(&items); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteJsonExit(items)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(
			client.ContentJson().PostContent(ctx, "/parse-tag-array", []map[string]any{
				{
					"title": "  Foo  ",
					"tags":  []string{" A ", " B "},
				},
				{
					"title": "  Bar  ",
					"tags":  []string{" C "},
				},
			}),
			`[{"title":"Foo","tags":["a","b"]},{"title":"Bar","tags":["c"]}]`,
		)
		t.Assert(
			client.ContentJson().PostContent(ctx, "/parse-tag-array", []map[string]any{
				{
					"title": "   ",
				},
			}),
			`The title field is required`,
		)
	})
}

func Benchmark_Params_ParseTag(b *testing.B) {
	type benchmarkNoParseReq struct {
		g.Meta `path:"/bench/no-parse" method:"post"`
		Title  string   `json:"title"`
		Email  string   `json:"email"`
		Tags   []string `json:"tags"`
	}
	type benchmarkWithParseReq struct {
		g.Meta `path:"/bench/with-parse" method:"post"`
		Title  string   `json:"title" parse:"trim-space|lower"`
		Email  string   `json:"email" parse:"trim-space|lower"`
		Tags   []string `json:"tags" parse:"foreach|trim-space|lower"`
	}
	type benchmarkArrayNoParseItem struct {
		Title string   `json:"title"`
		Tags  []string `json:"tags"`
	}
	type benchmarkArrayWithParseItem struct {
		Title string   `json:"title" parse:"trim-space|lower"`
		Tags  []string `json:"tags" parse:"foreach|trim-space|lower"`
	}

	b.StopTimer()

	s := g.Server(guid.S())
	s.SetDumpRouterMap(false)
	s.SetAccessLogEnabled(false)
	s.SetErrorLogEnabled(false)
	s.BindHandler("/bench/no-parse", func(ctx context.Context, req *benchmarkNoParseReq) (res any, err error) {
		g.RequestFromCtx(ctx).Response.Write(req.Title)
		return nil, nil
	})
	s.BindHandler("/bench/with-parse", func(ctx context.Context, req *benchmarkWithParseReq) (res any, err error) {
		g.RequestFromCtx(ctx).Response.Write(req.Title)
		return nil, nil
	})
	s.BindHandler("/bench/array-no-parse", func(r *ghttp.Request) {
		var items []*benchmarkArrayNoParseItem
		if err := r.Parse(&items); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.Write(items[0].Title)
	})
	s.BindHandler("/bench/array-with-parse", func(r *ghttp.Request) {
		var items []*benchmarkArrayWithParseItem
		if err := r.Parse(&items); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.Write(items[0].Title)
	})
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	client := g.Client()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	jsonClient := client.ContentJson()

	structPayload := g.Map{
		"title": "  Demo Title  ",
		"email": "  Demo@Example.COM  ",
		"tags":  []string{" Alpha ", " Beta "},
	}
	arrayPayload := []map[string]any{
		{
			"title": "  Demo Title  ",
			"tags":  []string{" Alpha ", " Beta "},
		},
		{
			"title": "  Extra  ",
			"tags":  []string{" Gamma "},
		},
	}

	b.StartTimer()
	b.Run("struct_no_parse_tag", func(b *testing.B) {
		b.ReportAllocs()
		jsonClient.PostContent(ctx, "/bench/no-parse", structPayload)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonClient.PostContent(ctx, "/bench/no-parse", structPayload)
		}
	})
	b.Run("struct_with_parse_tag", func(b *testing.B) {
		b.ReportAllocs()
		jsonClient.PostContent(ctx, "/bench/with-parse", structPayload)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonClient.PostContent(ctx, "/bench/with-parse", structPayload)
		}
	})
	b.Run("array_no_parse_tag", func(b *testing.B) {
		b.ReportAllocs()
		jsonClient.PostContent(ctx, "/bench/array-no-parse", arrayPayload)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonClient.PostContent(ctx, "/bench/array-no-parse", arrayPayload)
		}
	})
	b.Run("array_with_parse_tag", func(b *testing.B) {
		b.ReportAllocs()
		jsonClient.PostContent(ctx, "/bench/array-with-parse", arrayPayload)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			jsonClient.PostContent(ctx, "/bench/array-with-parse", arrayPayload)
		}
	})
}
