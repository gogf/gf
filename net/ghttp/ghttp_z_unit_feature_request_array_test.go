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
	"io"
	"mime/multipart"
	"net/url"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
// `type:"array"` tag, so JSON array request bodies are rejected.
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

// ArrayEmbeddedItems holds the array field in a tagged embedded request field.
type ArrayEmbeddedItems struct {
	Items []arrayItem `json:"items"`
}

// arrayReqEmbedded checks the runtime behavior of tagged embedded array fields.
type arrayReqEmbedded struct {
	g.Meta             `mime:"application/json" method:"post" path:"/array/embedded" type:"array"`
	ArrayEmbeddedItems `json:"payload"`
}

// arrayReqIgnored skips the hidden array field when choosing the body receiver.
type arrayReqIgnored struct {
	g.Meta `mime:"application/json" method:"post" path:"/array/ignored" type:"array"`
	Hidden []arrayItem `json:"-"`
	Items  []arrayItem `json:"items"`
}

// arrayReqNoReceiver declares an array body without a field that can receive it.
type arrayReqNoReceiver struct {
	g.Meta `type:"array"`
	Name   string `json:"name"`
}

// arrayRes is the response struct in the tests of the JSON array request body.
type arrayRes struct {
	ItemsCount   int    `json:"itemsCount"`
	ItemsNil     bool   `json:"itemsNil"`
	TagsCount    int    `json:"tagsCount"`
	FirstId      int    `json:"firstId"`
	FirstContent string `json:"firstContent"`
	ExtraField   string `json:"extraField"`
}

// startArrayServer starts a test server and waits for its listening port.
func startArrayServer(t *testing.T, s *ghttp.Server) {
	t.Helper()
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}
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
}

func writeArrayRes(r *ghttp.Request, items []arrayItem, tags []string, extraField string) {
	var res = arrayRes{
		ItemsCount: len(items),
		ItemsNil:   items == nil,
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
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(new(arrayController))
	})
	s.BindHandler("/array/raw", rawArrayHandler)
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

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
			expect: `{"itemsCount":2,"itemsNil":false,"tagsCount":0,"firstId":1,"firstContent":"content1","extraField":""}`,
		},
		{
			name:   "tagged request struct with single item",
			path:   "/array/tagged",
			body:   `[{"id":100,"name":"only one","content":"content100"}]`,
			expect: `{"itemsCount":1,"itemsNil":false,"tagsCount":0,"firstId":100,"firstContent":"content100","extraField":""}`,
		},
		{
			name:   "tagged request struct with empty array",
			path:   "/array/tagged",
			body:   `[]`,
			expect: `{"itemsCount":0,"itemsNil":false,"tagsCount":0,"firstId":0,"firstContent":"","extraField":""}`,
		},
		{
			// The JSON object request body is still working for the tagged request struct.
			name:   "tagged request struct with object body",
			path:   "/array/tagged",
			body:   `{"items":[{"id":7,"content":"content7"}]}`,
			expect: `{"itemsCount":1,"itemsNil":false,"tagsCount":0,"firstId":7,"firstContent":"content7","extraField":""}`,
		},
		{
			// It uses the request struct fields by reflection as the handler is not a strict route.
			name:   "none-strict route handler",
			path:   "/array/raw",
			body:   `[{"id":10,"content":"content10"},{"id":20,"content":"content20"}]`,
			expect: `{"itemsCount":2,"itemsNil":false,"tagsCount":0,"firstId":10,"firstContent":"content10","extraField":""}`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				t.Assert(client.PostContent(ctx, c.path, c.body), c.expect)
			})
		})
	}

	// The JSON array request body is not a valid request body for the request struct without
	// the `type:"array"` tag, which is rejected instead of being silently ignored with all
	// the attributes left as zero values.
	t.Run("untagged request struct is rejected", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			var (
				resp   = client.PostContent(ctx, "/array/untagged", `[{"id":1,"content":"content1"}]`)
				parsed = gjson.New(resp)
			)
			t.Assert(parsed.Get("code").Int(), gcode.CodeInvalidParameter.Code())
			t.Assert(parsed.Get("data").IsNil(), true)
		})
	})

	// The JSON array request body works with the JSON content type either.
	t.Run("tagged request struct with json content type", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(
				client.ContentJson().PostContent(ctx, "/array/tagged", `[{"id":3,"content":"content3"}]`),
				`{"itemsCount":1,"itemsNil":false,"tagsCount":0,"firstId":3,"firstContent":"content3","extraField":""}`,
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
	startArrayServer(t, s)

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

// arrayReqParams is the request struct used for checking that the JSON array request body
// is not misparsed as the form parameters.
type arrayReqParams struct {
	g.Meta `mime:"application/json" method:"post" path:"/array/params" type:"array"`
	Items  []arrayItem `json:"items"`
}

// Params returns the parsed array and request parameters for the fallback regression test.
func (c *arrayController) Params(ctx context.Context, req *arrayReqParams) (res *arrayRes, err error) {
	var r = g.RequestFromCtx(ctx)
	r.Response.WriteJson(g.Map{
		"itemsCount": len(req.Items),
		"params":     r.GetRequestMap(),
	})
	return
}

// Embedded reports whether the JSON array reached a tagged embedded field.
func (c *arrayController) Embedded(ctx context.Context, req *arrayReqEmbedded) (res *arrayRes, err error) {
	writeArrayRes(g.RequestFromCtx(ctx), req.Items, nil, "")
	return
}

// Ignored reports whether the JSON array reached the first visible array field.
func (c *arrayController) Ignored(ctx context.Context, req *arrayReqIgnored) (res *arrayRes, err error) {
	writeArrayRes(g.RequestFromCtx(ctx), req.Items, nil, "")
	return
}

// Test_Params_JsonArray_FieldSelection checks runtime conversion for ignored and
// tagged embedded array fields.
func Test_Params_JsonArray_FieldSelection(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(arrayController))
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	for _, path := range []string{"/array/embedded", "/array/ignored"} {
		t.Run(path, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				t.Assert(
					client.PostContent(ctx, path, `[{"id":7}]`),
					`{"itemsCount":1,"itemsNil":false,"tagsCount":0,"firstId":7,"firstContent":"","extraField":""}`,
				)
			})
		})
	}
}

// Test_Params_JsonArray_NoReceiver rejects a tagged array body with no eligible field.
func Test_Params_JsonArray_NoReceiver(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/array/no-receiver", func(r *ghttp.Request) {
		var req arrayReqNoReceiver
		err := r.Parse(&req)
		r.Response.WriteJson(g.Map{"invalid": gerror.Code(err) == gcode.CodeInvalidParameter})
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	client := g.Client().ContentJson()
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	for _, body := range []string{`[{"name":"x"}]`, `[]`} {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(client.PostContent(ctx, "/array/no-receiver", body), `{"invalid":true}`)
		})
	}
}

// Test_Params_JsonArray_RequestBodyParams checks that the JSON array request body is not
// misparsed as the query/form parameters by the default parameters decoding of parseBody.
//
// The request body contains `&` here, which used to be split into parts by `gstr.Parse`,
// and the parts that do not start with `[` were written into the request parameters, for
// example the `b=2"}]` part produced a bogus parameter `b`.
func Test_Params_JsonArray_RequestBodyParams(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(arrayController))
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	for _, contentJson := range []bool{false, true} {
		t.Run(fmt.Sprintf("json content type: %t", contentJson), func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				client := g.Client()
				client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
				if contentJson {
					client = client.ContentJson()
				} else {
					client = client.ContentType("text/plain")
				}
				var (
					body   = `[{"id":1,"name":"item1","content":"https://example.com/path?a=1&b=2"}]`
					parsed = gjson.New(client.PostContent(ctx, "/array/params", body))
				)
				t.Assert(parsed.Get("itemsCount").Int(), 1)
				t.Assert(len(parsed.Get("params").Map()), 0)
			})
		})
	}
}

// Test_Params_JsonArray_ReloadParam checks that reparsing a changed body does not reuse
// the previous object or array representation.
func Test_Params_JsonArray_ReloadParam(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/array/reload", func(r *ghttp.Request) {
		r.GetRequestMap()
		var replacement = r.GetQuery("replacement").String()
		r.Body = io.NopCloser(bytes.NewBufferString(replacement))
		r.ContentLength = int64(len(replacement))
		r.ReloadParam()

		var req arrayReqTagged
		if err := r.Parse(&req); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteJson(g.Map{
			"itemsCount": len(req.Items),
			"firstId":    req.Items[0].Id,
			"params":     r.GetRequestMap(),
		})
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	for _, testCase := range []struct {
		name        string
		original    string
		replacement string
		firstId     int
		paramCount  int
	}{
		{"array to object", `[{"id":1}]`, `{"items":[{"id":2}]}`, 2, 1},
		{"object to array", `{"items":[{"id":1}]}`, `[{"id":2}]`, 2, 0},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				client := g.Client().ContentJson()
				client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
				path := "/array/reload?replacement=" + url.QueryEscape(testCase.replacement)
				parsed := gjson.New(client.PostContent(ctx, path, testCase.original))
				t.Assert(parsed.Get("itemsCount").Int(), 1)
				t.Assert(parsed.Get("firstId").Int(), testCase.firstId)
				if testCase.paramCount == 0 {
					t.Assert(parsed.Get("params").Map(), g.Map{"replacement": testCase.replacement})
				}
				t.Assert(len(parsed.Get("params").Map()), testCase.paramCount+1)
			})
		})
	}
}

// Test_Params_ReloadParam_ChangedSources checks that reparsing discards form and query
// parameters that no longer exist in the modified request.
func Test_Params_ReloadParam_ChangedSources(t *testing.T) {
	makeMultipartBody := func(name string) (string, string) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		if err := writer.WriteField(name, "1"); err != nil {
			t.Fatal(err)
		}
		if err := writer.Close(); err != nil {
			t.Fatal(err)
		}
		return writer.FormDataContentType(), body.String()
	}
	oldMultipartType, oldMultipartBody := makeMultipartBody("old")
	newMultipartType, newMultipartBody := makeMultipartBody("new")

	s := g.Server(guid.S())
	s.BindHandler("/reload/form", func(r *ghttp.Request) {
		r.GetRequestMap()
		r.Body = io.NopCloser(bytes.NewBufferString("new=2"))
		r.ContentLength = int64(len("new=2"))
		r.ReloadParam()
		r.Response.WriteJson(r.GetRequestMap())
	})
	s.BindHandler("/reload/query", func(r *ghttp.Request) {
		r.GetRequestMap()
		r.URL.RawQuery = ""
		r.ReloadParam()
		r.Response.WriteJson(r.GetRequestMap())
	})
	s.BindHandler("/reload/multipart", func(r *ghttp.Request) {
		r.GetRequestMap()
		r.Body = io.NopCloser(bytes.NewBufferString(newMultipartBody))
		r.ContentLength = int64(len(newMultipartBody))
		r.Header.Set("Content-Type", newMultipartType)
		r.ReloadParam()
		r.Response.WriteJson(r.GetRequestMap())
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	client := g.Client().ContentType("application/x-www-form-urlencoded")
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	t.Run("changed form body", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(client.PostContent(ctx, "/reload/form", "old=1"), `{"new":"2"}`)
		})
	})
	t.Run("cleared query string", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(client.GetContent(ctx, "/reload/query?old=1"), `{}`)
		})
	})
	t.Run("changed multipart body", func(t *testing.T) {
		gtest.C(t, func(t *gtest.T) {
			multipartClient := g.Client().ContentType(oldMultipartType)
			multipartClient.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
			t.Assert(multipartClient.PostContent(ctx, "/reload/multipart", oldMultipartBody), `{"new":"1"}`)
		})
	})
}

// Test_Params_ReloadQuery_Multipart keeps parsed form values and uploaded files
// while refreshing the URL query string.
func Test_Params_ReloadQuery_Multipart(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("field", "kept"); err != nil {
		t.Fatal(err)
	}
	filePart, err := writer.CreateFormFile("upload", "sample.txt")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := filePart.Write([]byte("payload")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	s := g.Server(guid.S())
	s.SetFormParsingMemory(1)
	s.BindHandler("/reload/query-multipart", func(r *ghttp.Request) {
		r.GetRequestMap()
		r.URL.RawQuery = "new=2"
		r.ReloadQuery()

		params := r.GetRequestMap()
		uploads := r.GetUploadFiles("upload")
		if len(uploads) != 1 {
			r.Response.WriteJson(g.Map{"error": "upload missing"})
			return
		}
		file, err := uploads[0].Open()
		if err != nil {
			r.Response.WriteJson(g.Map{"error": err.Error()})
			return
		}
		content, readErr := io.ReadAll(file)
		closeErr := file.Close()
		if readErr != nil || closeErr != nil {
			r.Response.WriteJson(g.Map{"error": fmt.Sprint(readErr, closeErr)})
			return
		}
		r.Response.WriteJson(g.Map{
			"old":       params["old"],
			"new":       params["new"],
			"field":     params["field"],
			"file":      string(content),
			"formNew":   r.FormValue("new"),
			"formField": r.FormValue("field"),
		})
	})
	s.SetDumpRouterMap(false)
	startArrayServer(t, s)

	client := g.Client().ContentType(writer.FormDataContentType())
	client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
	gtest.C(t, func(t *gtest.T) {
		response := gjson.New(client.PostContent(ctx, "/reload/query-multipart?old=1", body.String()))
		t.Assert(response.Get("old").IsNil(), true)
		t.Assert(response.Get("new").String(), "2")
		t.Assert(response.Get("field").String(), "kept")
		t.Assert(response.Get("file").String(), "payload")
		t.Assert(response.Get("formNew").String(), "2")
		t.Assert(response.Get("formField").String(), "kept")
	})
}
