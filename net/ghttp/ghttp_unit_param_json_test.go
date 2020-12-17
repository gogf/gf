// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_Params_Json_Request(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Time  *time.Time
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"required|length:2,20|password3|same:password1#||密码强度不足|两次密码不一致"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.WriteExit(r.Get("id"), r.Get("name"))
	})
	s.BindHandler("/map", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			r.Response.WriteExit(m["id"], m["name"], m["password1"], m["password2"])
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
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/get", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}`), ``)
		t.Assert(client.GetContent("/map", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}`), ``)
		t.Assert(client.PostContent("/parse", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}`), `1john123Abc!@#123Abc!@#`)
		t.Assert(client.PostContent("/parse", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123"}`), `密码强度不足; 两次密码不一致`)
	})
}

func Test_Params_Json_Response(t *testing.T) {
	type User struct {
		Uid      int
		Name     string
		SiteUrl  string `json:"-"`
		NickName string `json:"nickname,omitempty"`
		Pass1    string `json:"password1"`
		Pass2    string `json:"password2"`
	}

	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/json1", func(r *ghttp.Request) {
		r.Response.WriteJson(User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		})
	})
	s.BindHandler("/json2", func(r *ghttp.Request) {
		r.Response.WriteJson(&User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		})
	})
	s.BindHandler("/json3", func(r *ghttp.Request) {
		type Message struct {
			Code  int    `json:"code"`
			Body  string `json:"body,omitempty"`
			Error string `json:"error,omitempty"`
		}
		type ResponseJson struct {
			Success  bool        `json:"success"`
			Data     interface{} `json:"data,omitempty"`
			ExtData  interface{} `json:"ext_data,omitempty"`
			Paginate interface{} `json:"paginate,omitempty"`
			Message  Message     `json:"message,omitempty"`
		}
		responseJson := &ResponseJson{
			Success: true,
			Data:    nil,
			ExtData: nil,
			Message: Message{3, "测试", "error"},
		}
		r.Response.WriteJson(responseJson)
	})
	s.BindHandler("/json4", func(r *ghttp.Request) {
		type Message struct {
			Code  int    `json:"code"`
			Body  string `json:"body,omitempty"`
			Error string `json:"error,omitempty"`
		}
		type ResponseJson struct {
			Success  bool        `json:"success"`
			Data     interface{} `json:"data,omitempty"`
			ExtData  interface{} `json:"ext_data,omitempty"`
			Paginate interface{} `json:"paginate,omitempty"`
			Message  *Message    `json:"message,omitempty"`
		}
		responseJson := ResponseJson{
			Success: true,
			Data:    nil,
			ExtData: nil,
			Message: &Message{3, "测试", "error"},
		}
		r.Response.WriteJson(responseJson)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		map1 := make(map[string]interface{})
		err1 := json.Unmarshal([]byte(client.GetContent("/json1")), &map1)
		t.Assert(err1, nil)
		t.Assert(len(map1), 4)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		map2 := make(map[string]interface{})
		err2 := json.Unmarshal([]byte(client.GetContent("/json2")), &map2)
		t.Assert(err2, nil)
		t.Assert(len(map2), 4)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["Uid"], 100)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")

		map3 := make(map[string]interface{})
		err3 := json.Unmarshal([]byte(client.GetContent("/json3")), &map3)
		t.Assert(err3, nil)
		t.Assert(len(map3), 2)
		t.Assert(map3["success"], "true")
		t.Assert(map3["message"], g.Map{"body": "测试", "code": 3, "error": "error"})

		map4 := make(map[string]interface{})
		err4 := json.Unmarshal([]byte(client.GetContent("/json4")), &map4)
		t.Assert(err4, nil)
		t.Assert(len(map4), 2)
		t.Assert(map4["success"], "true")
		t.Assert(map4["message"], g.Map{"body": "测试", "code": 3, "error": "error"})
	})
}
