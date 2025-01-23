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

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Router_Handler_Strict_WithObject(t *testing.T) {
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
		t.Assert(client.GetContent(ctx, "/test/error"), `{"code":50,"message":"error","data":{"Id":1,"Age":0,"Name":""}}`)
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

func (ControllerForHandlerWithObjectAndMeta1) Index(ctx context.Context, req *TestForHandlerWithObjectAndMeta1Req) (res *TestForHandlerWithObjectAndMeta1Res, err error) {
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

func Test_Router_Handler_Strict_WithObjectAndMeta(t *testing.T) {
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

		t.Assert(client.GetContent(ctx, "/"), `{"code":65,"message":"Not Found","data":null}`)
		t.Assert(client.GetContent(ctx, "/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/custom-test2?age=18&name=john"), `{"code":65,"message":"Not Found","data":null}`)
	})
}

func Test_Router_Handler_Strict_Group_Bind(t *testing.T) {
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

		t.Assert(client.GetContent(ctx, "/"), `{"code":65,"message":"Not Found","data":null}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.PostContent(ctx, "/api/v1/custom-test2?age=18&name=john"), `{"code":65,"message":"Not Found","data":null}`)

		t.Assert(client.GetContent(ctx, "/api/v1/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v1/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)

		t.Assert(client.GetContent(ctx, "/api/v2/custom-test1?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test2?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test3?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Age":18}}`)
		t.Assert(client.GetContent(ctx, "/api/v2/custom-test4?age=18&name=john"), `{"code":0,"message":"","data":{"Id":1,"Name":"john"}}`)
	})
}

func Test_Issue1708(t *testing.T) {
	type Test struct {
		Name string `json:"name"`
	}
	type Req struct {
		Page       int      `json:"page"       dc:"分页码"`
		Size       int      `json:"size"       dc:"分页数量"`
		TargetType string   `json:"targetType" v:"required#评论内容类型错误" dc:"评论类型: topic/ask/article/reply"`
		TargetId   uint     `json:"targetId"   v:"required#评论目标ID错误" dc:"对应内容ID"`
		Test       [][]Test `json:"test"`
	}
	type Res struct {
		Page       int      `json:"page"       dc:"分页码"`
		Size       int      `json:"size"       dc:"分页数量"`
		TargetType string   `json:"targetType" v:"required#评论内容类型错误" dc:"评论类型: topic/ask/article/reply"`
		TargetId   uint     `json:"targetId"   v:"required#评论目标ID错误" dc:"对应内容ID"`
		Test       [][]Test `json:"test"`
	}

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test", func(ctx context.Context, req *Req) (res *Res, err error) {
		return &Res{
			Page:       req.Page,
			Size:       req.Size,
			TargetType: req.TargetType,
			TargetId:   req.TargetId,
			Test:       req.Test,
		}, nil
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		content := `
{
    "targetType":"topic",
    "targetId":10785,
    "test":[
        [
            {
                "name":"123"
            }
        ]
    ]
}
`
		t.Assert(
			client.PostContent(ctx, "/test", content),
			`{"code":0,"message":"","data":{"page":0,"size":0,"targetType":"topic","targetId":10785,"test":[[{"name":"123"}]]}}`,
		)
	})
}

func Test_Custom_Slice_Type_Attribute(t *testing.T) {
	type (
		WhiteListKey    string
		WhiteListValues []string
		WhiteList       map[WhiteListKey]WhiteListValues
	)
	type Req struct {
		Id   int
		List WhiteList
	}
	type Res struct {
		Content string
	}

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test", func(ctx context.Context, req *Req) (res *Res, err error) {
		return &Res{
			Content: gjson.MustEncodeString(req),
		}, nil
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		content := `
{
    "id":1,
	"list":{
		"key": ["a", "b", "c"]
	}
}
`
		t.Assert(
			client.PostContent(ctx, "/test", content),
			`{"code":0,"message":"","data":{"Content":"{\"Id\":1,\"List\":{\"key\":[\"a\",\"b\",\"c\"]}}"}}`,
		)
	})
}

func Test_Router_Handler_Strict_WithGeneric(t *testing.T) {
	type TestReq struct {
		Age int
	}
	type TestGeneric[T any] struct {
		Test T
	}
	type Test1Res struct {
		Age TestGeneric[int]
	}
	type Test2Res TestGeneric[int]
	type TestGenericRes[T any] struct {
		Test T
	}

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/test1", func(ctx context.Context, req *TestReq) (res *Test1Res, err error) {
		return &Test1Res{
			Age: TestGeneric[int]{
				Test: req.Age,
			},
		}, nil
	})
	s.BindHandler("/test1_slice", func(ctx context.Context, req *TestReq) (res []Test1Res, err error) {
		return []Test1Res{
			Test1Res{
				Age: TestGeneric[int]{
					Test: req.Age,
				},
			},
		}, nil
	})
	s.BindHandler("/test2", func(ctx context.Context, req *TestReq) (res *Test2Res, err error) {
		return &Test2Res{
			Test: req.Age,
		}, nil
	})

	s.BindHandler("/test2_slice", func(ctx context.Context, req *TestReq) (res []Test2Res, err error) {
		return []Test2Res{
			Test2Res{
				Test: req.Age,
			},
		}, nil
	})

	s.BindHandler("/test3", func(ctx context.Context, req *TestReq) (res *TestGenericRes[int], err error) {
		return &TestGenericRes[int]{
			Test: req.Age,
		}, nil
	})

	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	s.BindHandler("/test3_slice", func(ctx context.Context, req *TestReq) (res []TestGenericRes[int], err error) {
		return []TestGenericRes[int]{
			TestGenericRes[int]{
				Test: req.Age,
			},
		}, nil
	})

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/test1?age=1"), `{"code":0,"message":"","data":{"Age":{"Test":1}}}`)
		t.Assert(client.GetContent(ctx, "/test1_slice?age=1"), `{"code":0,"message":"","data":[{"Age":{"Test":1}}]}`)
		t.Assert(client.GetContent(ctx, "/test2?age=2"), `{"code":0,"message":"","data":{"Test":2}}`)
		t.Assert(client.GetContent(ctx, "/test2_slice?age=2"), `{"code":0,"message":"","data":[{"Test":2}]}`)
		t.Assert(client.GetContent(ctx, "/test3?age=3"), `{"code":0,"message":"","data":{"Test":3}}`)
		t.Assert(client.GetContent(ctx, "/test3_slice?age=3"), `{"code":0,"message":"","data":[{"Test":3}]}`)
	})
}

type ParameterCaseSensitiveController struct{}

type ParameterCaseSensitiveControllerPathReq struct {
	g.Meta `path:"/api/*path" method:"post"`
	Path   string
}

type ParameterCaseSensitiveControllerPathRes struct {
	Path string
}

func (c *ParameterCaseSensitiveController) Path(
	ctx context.Context,
	req *ParameterCaseSensitiveControllerPathReq,
) (res *ParameterCaseSensitiveControllerPathRes, err error) {
	return &ParameterCaseSensitiveControllerPathRes{Path: req.Path}, nil
}

func Test_Router_Handler_Strict_ParameterCaseSensitive(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(&ParameterCaseSensitiveController{})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		for i := 0; i < 1000; i++ {
			t.Assert(
				client.PostContent(ctx, "/api/111", `{"Path":"222"}`),
				`{"code":0,"message":"","data":{"Path":"222"}}`,
			)
		}
	})
}

type testJsonRawMessageIssue3449Req struct {
	g.Meta `path:"/test" method:"POST" sm:"hello" tags:"示例"`

	Name    string          `json:"name" v:"required" dc:"名称"`
	JSONRaw json.RawMessage `json:"jsonRaw" dc:"原始JSON"`
}
type testJsonRawMessageIssue3449Res struct {
	Name    string          `json:"name" v:"required" dc:"名称"`
	JSONRaw json.RawMessage `json:"jsonRaw" dc:"原始JSON"`
}

type testJsonRawMessageIssue3449 struct {
}

func (t *testJsonRawMessageIssue3449) Test(ctx context.Context, req *testJsonRawMessageIssue3449Req) (res *testJsonRawMessageIssue3449Res, err error) {
	return &testJsonRawMessageIssue3449Res{
		Name:    req.Name,
		JSONRaw: req.JSONRaw,
	}, nil
}

// https://github.com/gogf/gf/issues/3449
func Test_JsonRawMessage_Issue3449(t *testing.T) {

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(testJsonRawMessageIssue3449))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		v1 := map[string]any{
			"jkey1": "11",
			"jkey2": "12",
		}

		v2 := map[string]any{
			"jkey1": "21",
			"jkey2": "22",
		}
		data := map[string]any{
			"Name": "test",
			"jsonRaw": []any{
				v1, v2,
			},
		}

		expect1 := `{"code":0,"message":"","data":{"name":"test","jsonRaw":[{"jkey1":"11","jkey2":"12"},{"jkey1":"21","jkey2":"22"}]}}`
		t.Assert(client.PostContent(ctx, "/test", data), expect1)

		expect2 := `{"code":0,"message":"","data":{"name":"test","jsonRaw":{"jkey1":"11","jkey2":"12"}}}`
		t.Assert(client.PostContent(ctx, "/test", map[string]any{
			"Name":    "test",
			"jsonRaw": v1,
		}), expect2)

	})
}

type testNullStringIssue3465Req struct {
	g.Meta `path:"/test" method:"get" sm:"hello" tags:"示例"`
	Name   []string `json:"name" v:"required"`
}
type testNullStringIssue3465Res struct {
	Name []string `json:"name" v:"required" `
}

type testNullStringIssue3465 struct {
}

func (t *testNullStringIssue3465) Test(ctx context.Context, req *testNullStringIssue3465Req) (res *testNullStringIssue3465Res, err error) {
	return &testNullStringIssue3465Res{
		Name: req.Name,
	}, nil
}

// https://github.com/gogf/gf/issues/3465
func Test_NullString_Issue3465(t *testing.T) {

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(new(testNullStringIssue3465))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		data1 := map[string]any{
			"name": "null",
		}

		expect1 := `{"code":0,"message":"","data":{"name":["null"]}}`
		t.Assert(client.GetContent(ctx, "/test", data1), expect1)

		data2 := map[string]any{
			"name": []string{"null", "null"},
		}
		expect2 := `{"code":0,"message":"","data":{"name":["null","null"]}}`
		t.Assert(client.GetContent(ctx, "/test", data2), expect2)

	})
}
