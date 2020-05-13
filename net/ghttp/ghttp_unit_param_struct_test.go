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

	"github.com/gogf/gf/util/gvalid"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_Params_Struct(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Time  *time.Time
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"passwd1 @required|length:2,20|password3#||密码强度不足"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/struct1", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := new(User)
			if err := r.GetStruct(user); err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(user.Id, user.Name, user.Pass1, user.Pass2)
		}
	})
	s.BindHandler("/struct2", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := (*User)(nil)
			if err := r.GetStruct(&user); err != nil {
				r.Response.WriteExit(err)
			}
			if user != nil {
				r.Response.WriteExit(user.Id, user.Name, user.Pass1, user.Pass2)
			}
		}
	})
	s.BindHandler("/struct-valid", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			user := new(User)
			if err := r.GetStruct(user); err != nil {
				r.Response.WriteExit(err)
			}
			if err := gvalid.CheckStruct(user, nil); err != nil {
				r.Response.WriteExit(err)
			}
		}
	})
	s.BindHandler("/parse", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			var user *User
			if err := r.Parse(&user); err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(user.Id, user.Name, user.Pass1, user.Pass2)
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent("/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent("/struct2", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent("/struct2", ``), ``)
		t.Assert(client.PostContent("/struct-valid", `id=1&name=john&password1=123&password2=0`), `The value length must be between 2 and 20; 密码强度不足`)
		t.Assert(client.PostContent("/parse", `id=1&name=john&password1=123&password2=0`), `The value length must be between 2 and 20; 密码强度不足`)
		t.Assert(client.GetContent("/parse", `id=1&name=john&password1=123&password2=456`), `密码强度不足`)
		t.Assert(client.GetContent("/parse", `id=1&name=john&password1=123Abc!@#&password2=123Abc!@#`), `1john123Abc!@#123Abc!@#`)
		t.Assert(client.PostContent("/parse", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}`), `1john123Abc!@#123Abc!@#`)
	})
}

func Test_Params_Structs(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Time  *time.Time
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"passwd1 @required|length:2,20|password3#||密码强度不足"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse1", func(r *ghttp.Request) {
		var users []*User
		if err := r.Parse(&users); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(users[0].Id, users[1].Id)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.PostContent(
			"/parse1",
			`[{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}, {"id":2,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}]`),
			`12`,
		)
	})
}
