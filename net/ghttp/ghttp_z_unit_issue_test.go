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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

// https://github.com/gogf/gf/issues/1609
func Test_Issue1609(t *testing.T) {
	s := g.Server(guid.S())
	group := s.Group("/api/get")
	group.GET("/", func(r *ghttp.Request) {
		r.Response.Write("get")
	})
	s.SetDumpRouterMap(false)
	gtest.Assert(s.Start(), nil)
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/api/get"), "get")
		t.Assert(c.PostContent(ctx, "/test"), "Not Found")
	})
}

func Test_Issue1611(t *testing.T) {
	s := g.Server(guid.S())
	v := g.View(guid.S())
	content := "This is header"
	gtest.AssertNil(v.SetPath(gtest.DataPath("issue1611")))
	s.SetView(v)
	s.BindHandler("/", func(r *ghttp.Request) {
		gtest.AssertNil(r.Response.WriteTpl("index/layout.html", g.Map{
			"header": content,
		}))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(gstr.Contains(c.GetContent(ctx, "/"), content), true)
	})
}

// https://github.com/gogf/gf/issues/1626
func Test_Issue1626(t *testing.T) {
	type TestReq struct {
		Name string `v:"required"`
	}
	type TestRes struct {
		Name string
	}
	s := g.Server(guid.S())
	s.Use(
		ghttp.MiddlewareHandlerResponse,
		func(r *ghttp.Request) {
			r.Middleware.Next()
			if err := r.GetError(); err != nil {
				r.Response.ClearBuffer()
				r.Response.Write(err.Error())
			}
		},
	)
	s.BindHandler("/test", func(ctx context.Context, req *TestReq) (res *TestRes, err error) {
		return &TestRes{Name: req.Name}, nil
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/test"), `The Name field is required`)
		t.Assert(
			gstr.Contains(c.GetContent(ctx, "/test?name=john"), `{"Name":"john"}`),
			true,
		)
	})
}

type Issue1653TestReq struct {
	g.Meta    `path:"/test" method:"post" summary:"执行报表查询" tags:""`
	UUID      string  `json:"uuid" v:"required#菜单唯一码不可为空" dc:""`
	Limit     int     `json:"limit"`
	Filter    []g.Map `json:"filter"`
	FilterMap g.Map   `json:"filter_map"`
}

type Issue1653TestRes struct {
	UUID     string      `json:"uuid"`
	FeedBack interface{} `json:"feed_back"`
}

type cIssue1653Foo struct{}

var Issue1653Foo = new(cIssue1653Foo)

func (r cIssue1653Foo) PostTest(ctx context.Context, req *Issue1653TestReq) (*Issue1653TestRes, error) {
	return &Issue1653TestRes{UUID: req.UUID, FeedBack: req.Filter[0]["code"]}, nil
}

func Test_Issue1653(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/boot", func(grp *ghttp.RouterGroup) {
		grp.Bind(Issue1653Foo)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(1000 * time.Millisecond)
	// g.Client()测试：
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		dataReq := `
{"uuid":"28ee701c-7daf-4cdc-9a62-6d6704e6112b","limit":0,"filter":
[
{
"code":"P00001","constraint":"",
"created_at":"2022-03-08 04:56:15","created_by":"3ed72aba-1622-4262-a61e-83581e020763","default_value":"MonthStart()",
"expression":"AND A.DLVDAT_0>='%v'","force":false,"frequent":true,"name":"发货日期起",
"parent":"13109602-0da3-49b9-827f-2f44183ab756","read_only":false,"reference":null,"type":"date",
"updated_at":"2022-03-08 04:56:15","updated_by":"3ed72aba-1622-4262-a61e-83581e020763","updated_tick":1,
"uuid":"e6cd3268-1d75-42e0-83f9-f1f7b29976e8"
},
{
"code":"P00002","constraint":"","created_at":"2022-03-08 04:56:15","created_by":
"3ed72aba-1622-4262-a61e-83581e020763","default_value":"MonthEnd()","expression":"AND A.DLVDAT_0<='%v'","force":false,"frequent":true,
"name":"发货日期止","parent":"13109602-0da3-49b9-827f-2f44183ab756","read_only":false,"reference":null,"type":"date","updated_at":
"2022-03-08 04:56:15","updated_by":"3ed72aba-1622-4262-a61e-83581e020763","updated_tick":1,"uuid":"dba005b5-655e-4ac4-8b22-898aa3ad2294"
}
],
"filter_map":{"P00001":1646064000000,"P00002":1648742399999},
"selector_template":""
}
`
		resContent := c.PostContent(ctx, "/boot/test", dataReq)
		t.Assert(resContent, `{"code":0,"message":"","data":{"uuid":"28ee701c-7daf-4cdc-9a62-6d6704e6112b","feed_back":"P00001"}}`)
	})
}

type LbseMasterHead struct {
	Code     string   `json:"code" v:"code@required|min-length:1#The code is required"`
	Active   bool     `json:"active"`
	Preset   bool     `json:"preset"`
	Superior string   `json:"superior"`
	Path     []string `json:"path"`
	Sort     int      `json:"sort"`
	Folder   bool     `json:"folder"`
	Test     string   `json:"test" v:"required"`
}

type Template struct {
	LbseMasterHead
	Datasource string `json:"datasource" v:"required|length:32,32#The datasource is required"`
	SQLText    string `json:"sql_text"`
}

type TemplateCreateReq struct {
	g.Meta `path:"/test" method:"post" summary:"Create template" tags:"Template"`
	Master Template `json:"master"`
}

type TemplateCreateRes struct{}

type cFoo1 struct{}

var Foo1 = new(cFoo1)

func (r cFoo1) PostTest1(ctx context.Context, req *TemplateCreateReq) (res *TemplateCreateRes, err error) {
	g.Dump(req)
	return
}

// https://github.com/gogf/gf/issues/1662
func Test_Issue662(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/boot", func(grp *ghttp.RouterGroup) {
		grp.Bind(Foo1)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(1000 * time.Millisecond)

	// g.Client()测试：
	// code字段传入空字符串时，校验没有提示
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		dataReq := `
{"master":{"active":true,"code":"","created_at":"","created_by":"","created_by_text":"","datasource":"38b6f170-a584-43fc-8912-cc1e9bf1b1a9","description":"币种","folder":false,"path":"[\"XCUR\"]","preset":false,"sort":1000,"sql_text":"SELECT!!!!","superior":null,"updated_at":"","updated_by":"","updated_by_text":"","updated_tick":0,"uuid":""},"translation":[{"code":"zh_CN","text":"币种"},{"code":"en_US","text":"币种"}],"filters":null,"fields":[{"code":"F001","created_at":"2022-01-18 23:37:38","created_by":"3ed72aba-1622-4262-a61e-83581e020763","field":"value","hide":false,"min_width":120,"name":"value","parent":"296154bf-b718-4e8f-8b70-efb969b831ec","updated_at":"2022-01-18 23:37:38","updated_by":"3ed72aba-1622-4262-a61e-83581e020763","updated_tick":1,"uuid":"f2140b7a-044c-41c3-b70e-852e6160b21b"},{"code":"F002","created_at":"2022-01-18 23:37:38","created_by":"3ed72aba-1622-4262-a61e-83581e020763","field":"label","hide":false,"min_width":120,"name":"label","parent":"296154bf-b718-4e8f-8b70-efb969b831ec","updated_at":"2022-01-18 23:37:38","updated_by":"3ed72aba-1622-4262-a61e-83581e020763","updated_tick":1,"uuid":"2d3bba5d-308b-4dba-bcac-f093e6556eca"}],"limit":0}
`
		t.Assert(c.PostContent(ctx, "/boot/test", dataReq), `{"code":51,"message":"The code is required","data":null}`)
	})
}

type DemoReq struct {
	g.Meta `path:"/demo" method:"post"`
	Data   *gjson.Json
}

type DemoRes struct {
	Content string
}

type Api struct{}

func (a *Api) Demo(ctx context.Context, req *DemoReq) (res *DemoRes, err error) {
	return &DemoRes{
		Content: req.Data.MustToJsonString(),
	}, err
}

var api = Api{}

// https://github.com/gogf/gf/issues/2172
func Test_Issue2172(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(api)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(1000 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		dataReq := `{"data":{"asd":1}}`
		t.Assert(c.PostContent(ctx, "/demo", dataReq), `{"code":0,"message":"","data":{"Content":"{\"asd\":1}"}}`)
	})
}

// https://github.com/gogf/gf/issues/2334
func Test_Issue2334(t *testing.T) {
	s := g.Server(guid.S())
	s.SetServerRoot(gtest.DataPath("static1"))
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(1000 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(c.GetContent(ctx, "/index.html"), "index")

		c.SetHeader("If-Modified-Since", "Mon, 12 Dec 2040 05:53:35 GMT")
		res, _ := c.Get(ctx, "/index.html")
		t.Assert(res.StatusCode, 304)
	})
}
