// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package httpclient_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/sdk/httpclient/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_HttpClient_With_Default_Handler(t *testing.T) {
	type Req struct {
		g.Meta `path:"/get" method:"get"`
	}
	type Res struct {
		Uid  int
		Name string
	}

	s := g.Server(guid.S())
	s.BindHandler("/get", func(r *ghttp.Request) {
		res := ghttp.DefaultHandlerResponse{
			Data: Res{
				Uid:  1,
				Name: "test",
			},
		}
		r.Response.WriteJson(res)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := httpclient.New(httpclient.Config{
			URL: fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()),
		})
		var (
			req = &Req{}
			res = &Res{}
		)
		err := client.Request(gctx.New(), req, res)
		t.AssertNil(err)
		t.AssertEQ(res.Uid, 1)
		t.AssertEQ(res.Name, "test")
	})
}

type CustomHandler struct{}

func (c CustomHandler) HandleResponse(ctx context.Context, res *gclient.Response, out interface{}) error {
	defer res.Close()
	if pointer, ok := out.(*string); ok {
		*pointer = res.ReadAllString()
	} else {
		return gerror.NewCodef(gcode.CodeInvalidParameter, "[CustomHandler] expectedType:'*string', but realType:'%T'", out)
	}
	return nil
}

func Test_HttpClient_With_Custom_Handler(t *testing.T) {
	type Req struct {
		g.Meta `path:"/get" method:"get"`
	}

	s := g.Server(guid.S())
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.WriteExit("It is a test.")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	client := httpclient.New(httpclient.Config{
		URL:     fmt.Sprintf("127.0.0.1:%d", s.GetListenedPort()),
		Handler: CustomHandler{},
	})
	req := &Req{}
	gtest.C(t, func(t *gtest.T) {
		var res = new(string)
		err := client.Request(gctx.New(), req, res)
		t.AssertNil(err)
		t.AssertEQ(*res, "It is a test.")
	})
	gtest.C(t, func(t *gtest.T) {
		var res string
		err := client.Request(gctx.New(), req, res)
		t.AssertEQ(err, gerror.NewCodef(gcode.CodeInvalidParameter, "[CustomHandler] expectedType:'*string', but realType:'%T'", res))
	})
}
