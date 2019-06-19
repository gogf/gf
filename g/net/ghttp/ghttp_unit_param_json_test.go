// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func Test_Params_Json(t *testing.T) {
	type User struct {
		Uid      int
		Name     string
		SiteUrl  string `gconv:"-"`
		NickName string `gconv:"nickname, omitempty"`
		Pass1    string `gconv:"password1"`
		Pass2    string `gconv:"password2"`
	}

	p := ports.PopRand()
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
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		map1 := make(map[string]interface{})
		err1 := json.Unmarshal([]byte(client.GetContent("/json1")), &map1)
		gtest.Assert(err1, nil)
		gtest.Assert(len(map1), 4)
		gtest.Assert(map1["Name"], "john")
		gtest.Assert(map1["Uid"], 100)
		gtest.Assert(map1["password1"], "123")
		gtest.Assert(map1["password2"], "456")

		map2 := make(map[string]interface{})
		err2 := json.Unmarshal([]byte(client.GetContent("/json2")), &map2)
		gtest.Assert(err2, nil)
		gtest.Assert(len(map2), 4)
		gtest.Assert(map2["Name"], "john")
		gtest.Assert(map2["Uid"], 100)
		gtest.Assert(map2["password1"], "123")
		gtest.Assert(map2["password2"], "456")

		map3 := make(map[string]interface{})
		err3 := json.Unmarshal([]byte(client.GetContent("/json3")), &map3)
		gtest.Assert(err3, nil)
		gtest.Assert(len(map3), 2)
		gtest.Assert(map3["success"], "true")
		gtest.Assert(map3["message"], g.Map{"body": "测试", "code": 3, "error": "error"})

		map4 := make(map[string]interface{})
		err4 := json.Unmarshal([]byte(client.GetContent("/json4")), &map4)
		gtest.Assert(err4, nil)
		gtest.Assert(len(map4), 2)
		gtest.Assert(map4["success"], "true")
		gtest.Assert(map4["message"], g.Map{"body": "测试", "code": 3, "error": "error"})
	})
}
