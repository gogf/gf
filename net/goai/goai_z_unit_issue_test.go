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
	"github.com/gogf/gf/v2/test/gtest"
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

func (Issue3664) RequiredTag(ctx context.Context, req *Issue3664RequiredTagReq) (res *Issue3664RequiredTagRes, err error) {
	res = &Issue3664RequiredTagRes{}
	return
}

// https://github.com/gogf/gf/issues/3664
func TestIssue3664(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server()
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
		apiContent := c.GetContent(ctx, "/api.json")
		j, err := gjson.LoadJson(apiContent)
		t.AssertNil(err)
		t.Assert(j.Get(`paths./default.post.requestBody.required`).String(), "")
		t.Assert(j.Get(`paths./required-tag.post.requestBody.required`).String(), "true")
	})
}
