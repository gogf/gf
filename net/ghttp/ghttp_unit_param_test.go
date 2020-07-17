// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_Params_Basic(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Pass1 string `params:"password1"`
		Pass2 string `params:"password2"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	// GET
	s.BindHandler("/get", func(r *ghttp.Request) {
		if r.GetQuery("array") != nil {
			r.Response.Write(r.GetQuery("array"))
		}
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
		if r.GetQuery("map") != nil {
			r.Response.Write(r.GetQueryMap()["map"].(map[string]interface{})["b"])
		}
		if r.GetQuery("a") != nil {
			r.Response.Write(r.GetQueryMapStrStr()["a"])
		}
	})
	// PUT
	s.BindHandler("/put", func(r *ghttp.Request) {
		if r.Get("array") != nil {
			r.Response.Write(r.Get("array"))
		}
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
		if r.Get("a") != nil {
			r.Response.Write(r.GetMapStrStr()["a"])
		}
	})
	// POST
	s.BindHandler("/post", func(r *ghttp.Request) {
		if r.GetPost("array") != nil {
			r.Response.Write(r.GetPost("array"))
		}
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
		if r.GetPost("a") != nil {
			r.Response.Write(r.GetPostMapStrStr()["a"])
		}
	})
	// DELETE
	s.BindHandler("/delete", func(r *ghttp.Request) {
		if r.Get("array") != nil {
			r.Response.Write(r.Get("array"))
		}
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
		if r.Get("a") != nil {
			r.Response.Write(r.GetMapStrStr()["a"])
		}
	})
	// PATCH
	s.BindHandler("/patch", func(r *ghttp.Request) {
		if r.Get("array") != nil {
			r.Response.Write(r.Get("array"))
		}
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
		if r.Get("a") != nil {
			r.Response.Write(r.GetMapStrStr()["a"])
		}
	})
	// Form
	s.BindHandler("/form", func(r *ghttp.Request) {
		if r.Get("array") != nil {
			r.Response.Write(r.GetForm("array"))
		}
		if r.Get("slice") != nil {
			r.Response.Write(r.GetForm("slice"))
		}
		if r.Get("bool") != nil {
			r.Response.Write(r.GetFormBool("bool"))
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.GetFormFloat32("float32"))
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.GetFormFloat64("float64"))
		}
		if r.Get("int") != nil {
			r.Response.Write(r.GetFormInt("int"))
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.GetFormUint("uint"))
		}
		if r.Get("string") != nil {
			r.Response.Write(r.GetFormString("string"))
		}
		if r.Get("map") != nil {
			r.Response.Write(r.GetFormMap()["map"].(map[string]interface{})["b"])
		}
		if r.Get("a") != nil {
			r.Response.Write(r.GetFormMapStrStr()["a"])
		}
	})
	s.BindHandler("/map", func(r *ghttp.Request) {
		if m := r.GetQueryMap(); len(m) > 0 {
			r.Response.Write(m["name"])
			return
		}
		if m := r.GetMap(); len(m) > 0 {
			r.Response.Write(m["name"])
			return
		}
	})
	s.BindHandler("/raw", func(r *ghttp.Request) {
		r.Response.Write(r.GetRaw())
	})
	s.BindHandler("/json", func(r *ghttp.Request) {
		j, err := r.GetJson()
		if err != nil {
			r.Response.Write(err)
			return
		}
		r.Response.Write(j.Get("name"))
	})
	s.BindHandler("/struct", func(r *ghttp.Request) {
		if m := r.GetQueryMap(); len(m) > 0 {
			user := new(User)
			r.GetQueryStruct(user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
			return
		}
		if m := r.GetMap(); len(m) > 0 {
			user := new(User)
			r.GetStruct(user)
			r.Response.Write(user.Id, user.Name, user.Pass1, user.Pass2)
			return
		}
	})
	s.BindHandler("/struct-with-nil", func(r *ghttp.Request) {
		user := (*User)(nil)
		err := r.GetStruct(&user)
		r.Response.Write(err)
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
		if m := r.GetMap(); len(m) > 0 {
			user1 := new(UserWithBase1)
			user2 := new(UserWithBase2)
			r.GetStruct(user1)
			r.GetStruct(user2)
			r.Response.Write(user1.Id, user1.Name, user1.Pass1, user1.Pass2)
			r.Response.Write(user2.Id, user2.Name, user2.Pass.Pass1, user2.Pass.Pass2)
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
		// GET
		t.Assert(client.GetContent("/get", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.GetContent("/get", "slice=1&slice=2"), `2`)
		t.Assert(client.GetContent("/get", "bool=1"), `true`)
		t.Assert(client.GetContent("/get", "bool=0"), `false`)
		t.Assert(client.GetContent("/get", "float32=0.11"), `0.11`)
		t.Assert(client.GetContent("/get", "float64=0.22"), `0.22`)
		t.Assert(client.GetContent("/get", "int=-10000"), `-10000`)
		t.Assert(client.GetContent("/get", "int=10000"), `10000`)
		t.Assert(client.GetContent("/get", "uint=10000"), `10000`)
		t.Assert(client.GetContent("/get", "uint=9"), `9`)
		t.Assert(client.GetContent("/get", "string=key"), `key`)
		t.Assert(client.GetContent("/get", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.GetContent("/get", "a=1&b=2"), `1`)

		// PUT
		t.Assert(client.PutContent("/put", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PutContent("/put", "slice=1&slice=2"), `2`)
		t.Assert(client.PutContent("/put", "bool=1"), `true`)
		t.Assert(client.PutContent("/put", "bool=0"), `false`)
		t.Assert(client.PutContent("/put", "float32=0.11"), `0.11`)
		t.Assert(client.PutContent("/put", "float64=0.22"), `0.22`)
		t.Assert(client.PutContent("/put", "int=-10000"), `-10000`)
		t.Assert(client.PutContent("/put", "int=10000"), `10000`)
		t.Assert(client.PutContent("/put", "uint=10000"), `10000`)
		t.Assert(client.PutContent("/put", "uint=9"), `9`)
		t.Assert(client.PutContent("/put", "string=key"), `key`)
		t.Assert(client.PutContent("/put", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PutContent("/put", "a=1&b=2"), `1`)

		// POST
		t.Assert(client.PostContent("/post", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PostContent("/post", "slice=1&slice=2"), `2`)
		t.Assert(client.PostContent("/post", "bool=1"), `true`)
		t.Assert(client.PostContent("/post", "bool=0"), `false`)
		t.Assert(client.PostContent("/post", "float32=0.11"), `0.11`)
		t.Assert(client.PostContent("/post", "float64=0.22"), `0.22`)
		t.Assert(client.PostContent("/post", "int=-10000"), `-10000`)
		t.Assert(client.PostContent("/post", "int=10000"), `10000`)
		t.Assert(client.PostContent("/post", "uint=10000"), `10000`)
		t.Assert(client.PostContent("/post", "uint=9"), `9`)
		t.Assert(client.PostContent("/post", "string=key"), `key`)
		t.Assert(client.PostContent("/post", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PostContent("/post", "a=1&b=2"), `1`)

		// DELETE
		t.Assert(client.DeleteContent("/delete", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.DeleteContent("/delete", "slice=1&slice=2"), `2`)
		t.Assert(client.DeleteContent("/delete", "bool=1"), `true`)
		t.Assert(client.DeleteContent("/delete", "bool=0"), `false`)
		t.Assert(client.DeleteContent("/delete", "float32=0.11"), `0.11`)
		t.Assert(client.DeleteContent("/delete", "float64=0.22"), `0.22`)
		t.Assert(client.DeleteContent("/delete", "int=-10000"), `-10000`)
		t.Assert(client.DeleteContent("/delete", "int=10000"), `10000`)
		t.Assert(client.DeleteContent("/delete", "uint=10000"), `10000`)
		t.Assert(client.DeleteContent("/delete", "uint=9"), `9`)
		t.Assert(client.DeleteContent("/delete", "string=key"), `key`)
		t.Assert(client.DeleteContent("/delete", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.DeleteContent("/delete", "a=1&b=2"), `1`)

		// PATCH
		t.Assert(client.PatchContent("/patch", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PatchContent("/patch", "slice=1&slice=2"), `2`)
		t.Assert(client.PatchContent("/patch", "bool=1"), `true`)
		t.Assert(client.PatchContent("/patch", "bool=0"), `false`)
		t.Assert(client.PatchContent("/patch", "float32=0.11"), `0.11`)
		t.Assert(client.PatchContent("/patch", "float64=0.22"), `0.22`)
		t.Assert(client.PatchContent("/patch", "int=-10000"), `-10000`)
		t.Assert(client.PatchContent("/patch", "int=10000"), `10000`)
		t.Assert(client.PatchContent("/patch", "uint=10000"), `10000`)
		t.Assert(client.PatchContent("/patch", "uint=9"), `9`)
		t.Assert(client.PatchContent("/patch", "string=key"), `key`)
		t.Assert(client.PatchContent("/patch", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PatchContent("/patch", "a=1&b=2"), `1`)

		// Form
		t.Assert(client.PostContent("/form", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PostContent("/form", "slice=1&slice=2"), `2`)
		t.Assert(client.PostContent("/form", "bool=1"), `true`)
		t.Assert(client.PostContent("/form", "bool=0"), `false`)
		t.Assert(client.PostContent("/form", "float32=0.11"), `0.11`)
		t.Assert(client.PostContent("/form", "float64=0.22"), `0.22`)
		t.Assert(client.PostContent("/form", "int=-10000"), `-10000`)
		t.Assert(client.PostContent("/form", "int=10000"), `10000`)
		t.Assert(client.PostContent("/form", "uint=10000"), `10000`)
		t.Assert(client.PostContent("/form", "uint=9"), `9`)
		t.Assert(client.PostContent("/form", "string=key"), `key`)
		t.Assert(client.PostContent("/form", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PostContent("/form", "a=1&b=2"), `1`)

		// Map
		t.Assert(client.GetContent("/map", "id=1&name=john"), `john`)
		t.Assert(client.PostContent("/map", "id=1&name=john"), `john`)

		// Raw
		t.Assert(client.PutContent("/raw", "id=1&name=john"), `id=1&name=john`)

		// Json
		t.Assert(client.PostContent("/json", `{"id":1, "name":"john"}`), `john`)

		// Struct
		t.Assert(client.GetContent("/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent("/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent("/struct-with-nil", ``), ``)
		t.Assert(client.PostContent("/struct-with-base", `id=1&name=john&password1=123&password2=456`), "1john1234561john123456")
	})
}

func Test_Params_SupportChars(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/form-value", func(r *ghttp.Request) {
		r.Response.Write(r.GetQuery("test-value"))
	})
	s.BindHandler("/form-array", func(r *ghttp.Request) {
		r.Response.Write(r.GetQuery("test-array"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := ghttp.NewClient()
		client.SetPrefix(prefix)

		t.Assert(client.PostContent("/form-value", "test-value=100"), "100")
		t.Assert(client.PostContent("/form-array", "test-array[]=1&test-array[]=2"), `["1","2"]`)
	})
}

func Test_Params_Priority(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Write(r.GetQuery("a"))
	})
	s.BindHandler("/post", func(r *ghttp.Request) {
		r.Response.Write(r.GetPost("a"))
	})
	s.BindHandler("/form", func(r *ghttp.Request) {
		r.Response.Write(r.GetForm("a"))
	})
	s.BindHandler("/request", func(r *ghttp.Request) {
		r.Response.Write(r.Get("a"))
	})
	s.BindHandler("/request-map", func(r *ghttp.Request) {
		r.Response.Write(r.GetMap(g.Map{
			"a": 1,
			"b": 2,
		}))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := ghttp.NewClient()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent("/query?a=1", "a=100"), "1")
		t.Assert(client.PostContent("/post?a=1", "a=100"), "100")
		t.Assert(client.PostContent("/form?a=1", "a=100"), "100")
		t.Assert(client.PutContent("/form?a=1", "a=100"), "100")
		t.Assert(client.GetContent("/request?a=1", "a=100"), "100")
		t.Assert(client.GetContent("/request-map?a=1&b=2&c=3", "a=100&b=200&c=300"), `{"a":"100","b":"200"}`)
	})
}

func Test_Params_GetRequestMap(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/map", func(r *ghttp.Request) {
		r.Response.Write(r.GetRequestMap())
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := ghttp.NewClient()
		client.SetPrefix(prefix)

		t.Assert(
			client.PostContent(
				"/map",
				"time_end2020-04-18 16:11:58&returnmsg=Success&attach=",
			),
			`{"attach":"","returnmsg":"Success"}`,
		)
	})
}
