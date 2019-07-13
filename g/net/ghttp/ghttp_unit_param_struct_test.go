// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/g/util/gvalid"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Params_Struct(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Pass1 string `params:"password1"`
		Pass2 string `params:"password2" gvalid:"passwd1 @required|length:2,20|password3#||密码强度不足"`
	}
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/struct1", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := new(User)
			r.GetToStruct(user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
		}
	})
	s.BindHandler("/struct2", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := (*User)(nil)
			r.GetToStruct(&user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
		}
	})
	s.BindHandler("/struct-valid", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := new(User)
			r.GetToStruct(user)
			err := gvalid.CheckStruct(user, nil)
			r.Response.Write(err.Maps())
		}
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		gtest.Assert(client.PostContent("/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		gtest.Assert(client.PostContent("/struct2", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		gtest.Assert(client.PostContent("/struct-valid", `id=1&name=john&password1=123&password2=0`), `{"passwd1":{"length":"字段长度为2到20个字符","password3":"密码强度不足"}}`)
	})
}
