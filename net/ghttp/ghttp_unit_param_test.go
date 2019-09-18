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

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_Params_Basic(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Pass1 string `params:"password1"`
		Pass2 string `params:"password2"`
	}
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/get", func(r *ghttp.Request) {
		if r.GetQuery("slice") != nil {
			r.Response.Write(r.GetQuery("slice"))
		}
		if r.GetQuery("bool") != nil {
			r.Response.Write(r.GetQueryBool("bool"))
		}
		if r.GetQuery("float32") != nil {
			r.Response.Write(r.GetQueryFloat32("float32"))
		}
		if r.GetQuery("float64") != nil {
			r.Response.Write(r.GetQueryFloat64("float64"))
		}
		if r.GetQuery("int") != nil {
			r.Response.Write(r.GetQueryInt("int"))
		}
		if r.GetQuery("uint") != nil {
			r.Response.Write(r.GetQueryUint("uint"))
		}
		if r.GetQuery("string") != nil {
			r.Response.Write(r.GetQueryString("string"))
		}
	})
	s.BindHandler("/put", func(r *ghttp.Request) {
		if r.Get("slice") != nil {
			r.Response.Write(r.Get("slice"))
		}
		if r.Get("bool") != nil {
			r.Response.Write(r.GetBool("bool"))
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.GetFloat32("float32"))
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.GetFloat64("float64"))
		}
		if r.Get("int") != nil {
			r.Response.Write(r.GetInt("int"))
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.GetUint("uint"))
		}
		if r.Get("string") != nil {
			r.Response.Write(r.GetString("string"))
		}
		if r.Get("map") != nil {
			r.Response.Write(r.GetMap()["map"].(map[string]interface{})["b"])
		}
	})
	s.BindHandler("/post", func(r *ghttp.Request) {
		if r.GetPost("slice") != nil {
			r.Response.Write(r.GetPost("slice"))
		}
		if r.GetPost("bool") != nil {
			r.Response.Write(r.GetPostBool("bool"))
		}
		if r.GetPost("float32") != nil {
			r.Response.Write(r.GetPostFloat32("float32"))
		}
		if r.GetPost("float64") != nil {
			r.Response.Write(r.GetPostFloat64("float64"))
		}
		if r.GetPost("int") != nil {
			r.Response.Write(r.GetPostInt("int"))
		}
		if r.GetPost("uint") != nil {
			r.Response.Write(r.GetPostUint("uint"))
		}
		if r.GetPost("string") != nil {
			r.Response.Write(r.GetPostString("string"))
		}
		if r.GetPost("map") != nil {
			r.Response.Write(r.GetPostMap()["map"].(map[string]interface{})["b"])
		}
	})
	s.BindHandler("/map", func(r *ghttp.Request) {
		if m := r.GetQueryMap(); len(m) > 0 {
			r.Response.Write(m["name"])
		}
		if m := r.GetPostMap(); len(m) > 0 {
			r.Response.Write(m["name"])
		}
	})
	s.BindHandler("/raw", func(r *ghttp.Request) {
		r.Response.Write(r.GetRaw())
	})
	s.BindHandler("/json", func(r *ghttp.Request) {
		r.Response.Write(r.GetJson().Get("name"))
	})
	s.BindHandler("/struct", func(r *ghttp.Request) {
		if m := r.GetQueryMap(); len(m) > 0 {
			user := new(User)
			r.GetQueryToStruct(user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
		}
		if m := r.GetPostMap(); len(m) > 0 {
			user := new(User)
			r.GetPostToStruct(user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
		}
	})
	s.BindHandler("/struct-with-base", func(r *ghttp.Request) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase1 struct {
			Id   int
			Name string
			Base
		}
		type UserWithBase2 struct {
			Id   int
			Name string
			Pass Base
		}
		if m := r.GetPostMap(); len(m) > 0 {
			user1 := new(UserWithBase1)
			user2 := new(UserWithBase2)
			r.GetToStruct(user1)
			r.GetToStruct(user2)
			r.Response.Write(user1.Id, user1.Name, user1.Pass1, user1.Pass2)
			r.Response.Write(user2.Id, user2.Name, user2.Pass.Pass1, user2.Pass.Pass2)
		}
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(200 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		// GET
		gtest.Assert(client.GetContent("/get", "slice=1&slice=2"), `2`)
		gtest.Assert(client.GetContent("/get", "bool=1"), `true`)
		gtest.Assert(client.GetContent("/get", "bool=0"), `false`)
		gtest.Assert(client.GetContent("/get", "float32=0.11"), `0.11`)
		gtest.Assert(client.GetContent("/get", "float64=0.22"), `0.22`)
		gtest.Assert(client.GetContent("/get", "int=-10000"), `-10000`)
		gtest.Assert(client.GetContent("/get", "int=10000"), `10000`)
		gtest.Assert(client.GetContent("/get", "uint=10000"), `10000`)
		gtest.Assert(client.GetContent("/get", "uint=9"), `9`)
		gtest.Assert(client.GetContent("/get", "string=key"), `key`)

		// PUT
		gtest.Assert(client.PutContent("/put", "slice=1&slice=2"), `2`)
		gtest.Assert(client.PutContent("/put", "bool=1"), `true`)
		gtest.Assert(client.PutContent("/put", "bool=0"), `false`)
		gtest.Assert(client.PutContent("/put", "float32=0.11"), `0.11`)
		gtest.Assert(client.PutContent("/put", "float64=0.22"), `0.22`)
		gtest.Assert(client.PutContent("/put", "int=-10000"), `-10000`)
		gtest.Assert(client.PutContent("/put", "int=10000"), `10000`)
		gtest.Assert(client.PutContent("/put", "uint=10000"), `10000`)
		gtest.Assert(client.PutContent("/put", "uint=9"), `9`)
		gtest.Assert(client.PutContent("/put", "string=key"), `key`)
		gtest.Assert(client.PutContent("/put", "map[a]=1&map[b]=2"), `2`)

		// POST
		gtest.Assert(client.PostContent("/post", "slice=1&slice=2"), `2`)
		gtest.Assert(client.PostContent("/post", "bool=1"), `true`)
		gtest.Assert(client.PostContent("/post", "bool=0"), `false`)
		gtest.Assert(client.PostContent("/post", "float32=0.11"), `0.11`)
		gtest.Assert(client.PostContent("/post", "float64=0.22"), `0.22`)
		gtest.Assert(client.PostContent("/post", "int=-10000"), `-10000`)
		gtest.Assert(client.PostContent("/post", "int=10000"), `10000`)
		gtest.Assert(client.PostContent("/post", "uint=10000"), `10000`)
		gtest.Assert(client.PostContent("/post", "uint=9"), `9`)
		gtest.Assert(client.PostContent("/post", "string=key"), `key`)
		gtest.Assert(client.PostContent("/post", "map[a]=1&map[b]=2"), `2`)

		// Map
		gtest.Assert(client.GetContent("/map", "id=1&name=john"), `john`)
		gtest.Assert(client.PostContent("/map", "id=1&name=john"), `john`)

		// Raw
		gtest.Assert(client.PutContent("/raw", "id=1&name=john"), `id=1&name=john`)

		// Json
		gtest.Assert(client.PostContent("/json", `{"id":1, "name":"john"}`), `john`)

		// Struct
		gtest.Assert(client.GetContent("/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		gtest.Assert(client.PostContent("/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		gtest.Assert(client.PostContent("/struct-with-base", `id=1&name=john&password1=123&password2=456`), "1john1234561john123456")
	})
}
