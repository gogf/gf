// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Client_DoRequestObj(t *testing.T) {
	type UserCreateReq struct {
		g.Meta `path:"/user" method:"post"`
		Id     int
		Name   string
	}
	type UserCreateRes struct {
		Id int
	}
	type UserQueryReq struct {
		g.Meta `path:"/user" method:"get"`
	}
	type UserQueryRes struct {
		Id   int
		Name string
	}
	p, _ := gtcp.GetFreePort()
	s := g.Server(p)
	s.Group("/user", func(group *ghttp.RouterGroup) {
		group.GET("/", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{"id": 1, "name": "john"})
		})
		group.POST("/", func(r *ghttp.Request) {
			r.Response.WriteJson(g.Map{"id": r.Get("Id")})
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client().SetPrefix(url).ContentJson()
		var (
			createRes *UserCreateRes
			createReq = UserCreateReq{
				Id:   1,
				Name: "john",
			}
		)
		err := client.DoRequestObj(ctx, createReq, &createRes)
		t.AssertNil(err)
		t.Assert(createRes.Id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client().SetPrefix(url).ContentJson()
		var (
			queryRes *UserQueryRes
			queryReq = UserQueryReq{}
		)
		err := client.DoRequestObj(ctx, queryReq, &queryRes)
		t.AssertNil(err)
		t.Assert(queryRes.Id, 1)
		t.Assert(queryRes.Name, "john")
	})
}
