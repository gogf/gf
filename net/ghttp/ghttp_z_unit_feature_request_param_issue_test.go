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

type Issue4087Req struct {
	g.Meta     `path:"/test" method:"post"`
	Page       int    `json:"page" example:"10" d:"1" v:"min:1#页码最小值不能低于1"  dc:"当前页码"`
	PerPage    int    `json:"pageSize" example:"1" d:"10" v:"min:1|max:200#每页数量最小值不能低于1|最大值不能大于200" dc:"每页数量"`
	Pagination bool   `json:"pagination" d:"true" dc:"是否需要进行分页"`
	Name       string `json:"name" d:"john"`
}

type Issue4087Res struct {
	g.Meta     `mime:"text/html" example:"string"`
	Page       int    `json:"page"`
	PerPage    int    `json:"pageSize"`
	Pagination bool   `json:"pagination"`
	Name       string `json:"name"`
}

type Issue4087GetReq struct {
	g.Meta     `path:"/test" method:"get"`
	Page       int    `json:"page" example:"10" d:"1" v:"min:1#页码最小值不能低于1"  dc:"当前页码"`
	PerPage    int    `json:"pageSize" example:"1" d:"10" v:"min:1|max:200#每页数量最小值不能低于1|最大值不能大于200" dc:"每页数量"`
	Pagination bool   `json:"pagination" d:"true" dc:"是否需要进行分页"`
	Name       string `json:"name" d:"john"`
}

type Issue4087GetRes struct {
	g.Meta     `mime:"text/html" example:"string"`
	Page       int    `json:"page"`
	PerPage    int    `json:"pageSize"`
	Pagination bool   `json:"pagination"`
	Name       string `json:"name"`
}

type Issue4087 struct{}

func (Issue4087) Default(ctx context.Context, req *Issue4087Req) (res *Issue4087Res, err error) {
	res = &Issue4087Res{
		Page:       req.Page,
		PerPage:    req.PerPage,
		Pagination: req.Pagination,
		Name:       req.Name,
	}
	return
}
func (Issue4087) DefaultGet(ctx context.Context, req *Issue4087GetReq) (res *Issue4087GetRes, err error) {
	res = &Issue4087GetRes{
		Page:       req.Page,
		PerPage:    req.PerPage,
		Pagination: req.Pagination,
		Name:       req.Name,
	}
	return
}

// https://github.com/gogf/gf/issues/4087
func Test_Issue4087(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareHandlerResponse)
		group.Bind(new(Issue4087))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		// post json
		t.Assert(client.ContentJson().PostContent(ctx, "/test", `{"pagination":true,"pageSize":1010,"name":"tom","page":10}`), `{"code":0,"message":"OK","data":{"page":10,"pageSize":1010,"pagination":true,"name":"tom"}}`)
		t.Assert(client.ContentJson().PostContent(ctx, "/test", `{"pagination":true,"pageSize":1010,"name":"tom","page":10,}`), `{"code":53,"message":"Parse Body failed: json.UnmarshalUseNumber failed: invalid character '}' looking for beginning of object key string","data":null}`)
		t.Assert(client.ContentJson().PostContent(ctx, "/test", ``), `{"code":0,"message":"OK","data":{"page":1,"pageSize":10,"pagination":true,"name":"john"}}`)

		// post xml
		t.Assert(client.ContentXml().PostContent(ctx, "/test", `<root><pagination>true</pagination><pageSize>1010</pageSize><name>tom</name><page>10</page></root>`), `{"code":0,"message":"OK","data":{"page":10,"pageSize":1010,"pagination":true,"name":"tom"}}`)
		t.Assert(client.ContentXml().PostContent(ctx, "/test", `<root><pagination>true</pagination><pageSize>1010</pageSize><name>tom</name><page>10</root>`), `{"code":53,"message":"Parse Body failed: mxj.NewMapXml failed: xml.Decoder.Token() - XML syntax error on line 1: element \u003cpage\u003e closed by \u003c/root\u003e","data":null}`)
		t.Assert(client.ContentXml().PostContent(ctx, "/test", ``), `{"code":0,"message":"OK","data":{"page":1,"pageSize":10,"pagination":true,"name":"john"}}`)

		// query
		t.Assert(client.GetContent(ctx, "/test", `pagination=true&pageSize=1010&name=tom&page=10`), `{"code":0,"message":"OK","data":{"page":10,"pageSize":1010,"pagination":true,"name":"tom"}}`)
		t.Assert(client.GetContent(ctx, "/test", `pagination=true&pageSize=1010&name=tom&page=10&pagination[]=true`), `{"code":53,"message":"Parse Query failed: expected type '[]any' for key 'pagination', but got 'string'","data":null}`)
		t.Assert(client.GetContent(ctx, "/test", ``), `{"code":0,"message":"OK","data":{"page":1,"pageSize":10,"pagination":true,"name":"john"}}`)

		// form
		mimePostForm := "application/x-www-form-urlencoded"
		t.Assert(client.SetContentType(mimePostForm).PostContent(ctx, "/test", `pagination=true&pageSize=1010&name=tom&page=10`), `{"code":0,"message":"OK","data":{"page":10,"pageSize":1010,"pagination":true,"name":"tom"}}`)
		t.Assert(client.SetContentType(mimePostForm).PostContent(ctx, "/test", `pagination=true&pageSize=1010&name=tom;&page=10&aaa`), `{"code":66,"message":"r.Request.ParseForm failed: invalid semicolon separator in query","data":null}`)
		t.Assert(client.SetContentType(mimePostForm).PostContent(ctx, "/test", ``), `{"code":0,"message":"OK","data":{"page":1,"pageSize":10,"pagination":true,"name":"john"}}`)

		t.Assert(client.SetContentType("multipart/form-data").PostContent(ctx, "/test", `{"pagination":true,"pageSize":1010,"name":"tom","page":10}`), `{"code":66,"message":"r.ParseMultipartForm failed: no multipart boundary param in Content-Type","data":null}`)
	})
}
