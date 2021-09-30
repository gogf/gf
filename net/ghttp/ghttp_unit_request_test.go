// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/gogf/gf/encoding/gjson"

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
			r.Response.Write(r.GetQuery("bool").Bool())
		}
		if r.GetQuery("float32") != nil {
			r.Response.Write(r.GetQuery("float32").Float32())
		}
		if r.GetQuery("float64") != nil {
			r.Response.Write(r.GetQuery("float64").Float64())
		}
		if r.GetQuery("int") != nil {
			r.Response.Write(r.GetQuery("int").Int())
		}
		if r.GetQuery("uint") != nil {
			r.Response.Write(r.GetQuery("uint").Uint())
		}
		if r.GetQuery("string") != nil {
			r.Response.Write(r.GetQuery("string").String())
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
			r.Response.Write(r.Get("bool").Bool())
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.Get("float32").Float32())
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.Get("float64").Float64())
		}
		if r.Get("int") != nil {
			r.Response.Write(r.Get("int").Int())
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.Get("uint").Uint())
		}
		if r.Get("string") != nil {
			r.Response.Write(r.Get("string").String())
		}
		if r.Get("map") != nil {
			r.Response.Write(r.GetMap()["map"].(map[string]interface{})["b"])
		}
		if r.Get("a") != nil {
			r.Response.Write(r.GetMapStrStr()["a"])
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
			r.Response.Write(r.Get("bool").Bool())
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.Get("float32").Float32())
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.Get("float64").Float64())
		}
		if r.Get("int") != nil {
			r.Response.Write(r.Get("int").Int())
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.Get("uint").Uint())
		}
		if r.Get("string") != nil {
			r.Response.Write(r.Get("string").String())
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
			r.Response.Write(r.Get("bool").Bool())
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.Get("float32").Float32())
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.Get("float64").Float64())
		}
		if r.Get("int") != nil {
			r.Response.Write(r.Get("int").Int())
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.Get("uint").Uint())
		}
		if r.Get("string") != nil {
			r.Response.Write(r.Get("string").String())
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
			r.Response.Write(r.GetForm("bool").Bool())
		}
		if r.Get("float32") != nil {
			r.Response.Write(r.GetForm("float32").Float32())
		}
		if r.Get("float64") != nil {
			r.Response.Write(r.GetForm("float64").Float64())
		}
		if r.Get("int") != nil {
			r.Response.Write(r.GetForm("int").Int())
		}
		if r.Get("uint") != nil {
			r.Response.Write(r.GetForm("uint").Uint())
		}
		if r.Get("string") != nil {
			r.Response.Write(r.GetForm("string").String())
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
		r.Response.Write(r.GetBody())
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
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		// GET
		t.Assert(client.GetContent(ctx, "/get", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.GetContent(ctx, "/get", "slice=1&slice=2"), `2`)
		t.Assert(client.GetContent(ctx, "/get", "bool=1"), `true`)
		t.Assert(client.GetContent(ctx, "/get", "bool=0"), `false`)
		t.Assert(client.GetContent(ctx, "/get", "float32=0.11"), `0.11`)
		t.Assert(client.GetContent(ctx, "/get", "float64=0.22"), `0.22`)
		t.Assert(client.GetContent(ctx, "/get", "int=-10000"), `-10000`)
		t.Assert(client.GetContent(ctx, "/get", "int=10000"), `10000`)
		t.Assert(client.GetContent(ctx, "/get", "uint=10000"), `10000`)
		t.Assert(client.GetContent(ctx, "/get", "uint=9"), `9`)
		t.Assert(client.GetContent(ctx, "/get", "string=key"), `key`)
		t.Assert(client.GetContent(ctx, "/get", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.GetContent(ctx, "/get", "a=1&b=2"), `1`)

		// PUT
		t.Assert(client.PutContent(ctx, "/put", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PutContent(ctx, "/put", "slice=1&slice=2"), `2`)
		t.Assert(client.PutContent(ctx, "/put", "bool=1"), `true`)
		t.Assert(client.PutContent(ctx, "/put", "bool=0"), `false`)
		t.Assert(client.PutContent(ctx, "/put", "float32=0.11"), `0.11`)
		t.Assert(client.PutContent(ctx, "/put", "float64=0.22"), `0.22`)
		t.Assert(client.PutContent(ctx, "/put", "int=-10000"), `-10000`)
		t.Assert(client.PutContent(ctx, "/put", "int=10000"), `10000`)
		t.Assert(client.PutContent(ctx, "/put", "uint=10000"), `10000`)
		t.Assert(client.PutContent(ctx, "/put", "uint=9"), `9`)
		t.Assert(client.PutContent(ctx, "/put", "string=key"), `key`)
		t.Assert(client.PutContent(ctx, "/put", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PutContent(ctx, "/put", "a=1&b=2"), `1`)

		// DELETE
		t.Assert(client.DeleteContent(ctx, "/delete", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.DeleteContent(ctx, "/delete", "slice=1&slice=2"), `2`)
		t.Assert(client.DeleteContent(ctx, "/delete", "bool=1"), `true`)
		t.Assert(client.DeleteContent(ctx, "/delete", "bool=0"), `false`)
		t.Assert(client.DeleteContent(ctx, "/delete", "float32=0.11"), `0.11`)
		t.Assert(client.DeleteContent(ctx, "/delete", "float64=0.22"), `0.22`)
		t.Assert(client.DeleteContent(ctx, "/delete", "int=-10000"), `-10000`)
		t.Assert(client.DeleteContent(ctx, "/delete", "int=10000"), `10000`)
		t.Assert(client.DeleteContent(ctx, "/delete", "uint=10000"), `10000`)
		t.Assert(client.DeleteContent(ctx, "/delete", "uint=9"), `9`)
		t.Assert(client.DeleteContent(ctx, "/delete", "string=key"), `key`)
		t.Assert(client.DeleteContent(ctx, "/delete", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.DeleteContent(ctx, "/delete", "a=1&b=2"), `1`)

		// PATCH
		t.Assert(client.PatchContent(ctx, "/patch", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PatchContent(ctx, "/patch", "slice=1&slice=2"), `2`)
		t.Assert(client.PatchContent(ctx, "/patch", "bool=1"), `true`)
		t.Assert(client.PatchContent(ctx, "/patch", "bool=0"), `false`)
		t.Assert(client.PatchContent(ctx, "/patch", "float32=0.11"), `0.11`)
		t.Assert(client.PatchContent(ctx, "/patch", "float64=0.22"), `0.22`)
		t.Assert(client.PatchContent(ctx, "/patch", "int=-10000"), `-10000`)
		t.Assert(client.PatchContent(ctx, "/patch", "int=10000"), `10000`)
		t.Assert(client.PatchContent(ctx, "/patch", "uint=10000"), `10000`)
		t.Assert(client.PatchContent(ctx, "/patch", "uint=9"), `9`)
		t.Assert(client.PatchContent(ctx, "/patch", "string=key"), `key`)
		t.Assert(client.PatchContent(ctx, "/patch", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PatchContent(ctx, "/patch", "a=1&b=2"), `1`)

		// Form
		t.Assert(client.PostContent(ctx, "/form", "array[]=1&array[]=2"), `["1","2"]`)
		t.Assert(client.PostContent(ctx, "/form", "slice=1&slice=2"), `2`)
		t.Assert(client.PostContent(ctx, "/form", "bool=1"), `true`)
		t.Assert(client.PostContent(ctx, "/form", "bool=0"), `false`)
		t.Assert(client.PostContent(ctx, "/form", "float32=0.11"), `0.11`)
		t.Assert(client.PostContent(ctx, "/form", "float64=0.22"), `0.22`)
		t.Assert(client.PostContent(ctx, "/form", "int=-10000"), `-10000`)
		t.Assert(client.PostContent(ctx, "/form", "int=10000"), `10000`)
		t.Assert(client.PostContent(ctx, "/form", "uint=10000"), `10000`)
		t.Assert(client.PostContent(ctx, "/form", "uint=9"), `9`)
		t.Assert(client.PostContent(ctx, "/form", "string=key"), `key`)
		t.Assert(client.PostContent(ctx, "/form", "map[a]=1&map[b]=2"), `2`)
		t.Assert(client.PostContent(ctx, "/form", "a=1&b=2"), `1`)

		// Map
		t.Assert(client.GetContent(ctx, "/map", "id=1&name=john"), `john`)
		t.Assert(client.PostContent(ctx, "/map", "id=1&name=john"), `john`)

		// Raw
		t.Assert(client.PutContent(ctx, "/raw", "id=1&name=john"), `id=1&name=john`)

		// Json
		t.Assert(client.PostContent(ctx, "/json", `{"id":1, "name":"john"}`), `john`)

		// Struct
		t.Assert(client.GetContent(ctx, "/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent(ctx, "/struct", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent(ctx, "/struct-with-nil", ``), ``)
		t.Assert(client.PostContent(ctx, "/struct-with-base", `id=1&name=john&password1=123&password2=456`), "1john1234561john")
	})
}

func Test_Params_Header(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/header", func(r *ghttp.Request) {
		r.Response.Write(r.GetHeader("test"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.Header(g.MapStrStr{"test": "123456"}).GetContent(ctx, "/header"), "123456")
	})
}

func Test_Params_SupportChars(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/form-value", func(r *ghttp.Request) {
		r.Response.Write(r.GetForm("test-value"))
	})
	s.BindHandler("/form-array", func(r *ghttp.Request) {
		r.Response.Write(r.GetForm("test-array"))
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.PostContent(ctx, "/form-value", "test-value=100"), "100")
		t.Assert(c.PostContent(ctx, "/form-array", "test-array[]=1&test-array[]=2"), `["1","2"]`)
	})
}

func Test_Params_Priority(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Write(r.GetQuery("a"))
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
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/query?a=1", "a=100"), "100")
		t.Assert(client.PostContent(ctx, "/form?a=1", "a=100"), "100")
		t.Assert(client.PutContent(ctx, "/form?a=1", "a=100"), "100")
		t.Assert(client.GetContent(ctx, "/request?a=1", "a=100"), "100")
		t.Assert(client.GetContent(ctx, "/request-map?a=1&b=2&c=3", "a=100&b=200&c=300"), `{"a":"100","b":"200"}`)
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
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(
			client.PostContent(ctx,
				"/map",
				"time_end2020-04-18 16:11:58&returnmsg=Success&attach=",
			),
			`{"attach":"","returnmsg":"Success"}`,
		)
	})
}

func Test_Params_Modify(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/param/modify", func(r *ghttp.Request) {
		param := r.GetMap()
		param["id"] = 2
		paramBytes, err := gjson.Encode(param)
		if err != nil {
			r.Response.Write(err)
			return
		}
		r.Request.Body = ioutil.NopCloser(bytes.NewReader(paramBytes))
		r.ReloadParam()
		r.Response.Write(r.GetMap())
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(
			client.PostContent(ctx,
				"/param/modify",
				`{"id":1}`,
			),
			`{"id":2}`,
		)
	})
}

func Test_Params_Parse_DefaultValueTag(t *testing.T) {
	type T struct {
		Name  string  `d:"john"`
		Score float32 `d:"60"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var t *T
		if err := r.Parse(&t); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(t)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.PostContent(ctx, "/parse"), `{"Name":"john","Score":60}`)
		t.Assert(client.PostContent(ctx, "/parse", `{"name":"smith"}`), `{"Name":"smith","Score":60}`)
		t.Assert(client.PostContent(ctx, "/parse", `{"name":"smith", "score":100}`), `{"Name":"smith","Score":100}`)
	})
}

func Test_Params_Parse_Validation(t *testing.T) {
	type RegisterReq struct {
		Name  string `p:"username"  v:"required|length:6,30#请输入账号|账号长度为:min到:max位"`
		Pass  string `p:"password1" v:"required|length:6,30#请输入密码|密码长度不够"`
		Pass2 string `p:"password2" v:"required|length:6,30|same:password1#请确认密码|密码长度不够|两次密码不一致"`
	}

	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var req *RegisterReq
		if err := r.Parse(&req); err != nil {
			r.Response.Write(err)
		} else {
			r.Response.Write("ok")
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/parse"), `请输入账号; 账号长度为6到30位; 请输入密码; 密码长度不够; 请确认密码; 密码长度不够; 两次密码不一致`)
		t.Assert(client.GetContent(ctx, "/parse?name=john11&password1=123456&password2=123"), `密码长度不够; 两次密码不一致`)
		t.Assert(client.GetContent(ctx, "/parse?name=john&password1=123456&password2=123456"), `账号长度为6到30位`)
		t.Assert(client.GetContent(ctx, "/parse?name=john11&password1=123456&password2=123456"), `ok`)
	})
}

func Test_Params_Parse_EmbeddedWithAliasName1(t *testing.T) {
	// 获取内容列表
	type ContentGetListInput struct {
		Type       string
		CategoryId uint
		Page       int
		Size       int
		Sort       int
		UserId     uint
	}
	// 获取内容列表
	type ContentGetListReq struct {
		ContentGetListInput
		CategoryId uint `p:"cate"`
		Page       int  `d:"1"  v:"min:0#分页号码错误"`
		Size       int  `d:"10" v:"max:50#分页数量最大50条"`
	}

	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var req *ContentGetListReq
		if err := r.Parse(&req); err != nil {
			r.Response.Write(err)
		} else {
			r.Response.Write(req.ContentGetListInput)
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/parse?cate=1&page=2&size=10"), `{"Type":"","CategoryId":0,"Page":2,"Size":10,"Sort":0,"UserId":0}`)
	})
}

func Test_Params_Parse_EmbeddedWithAliasName2(t *testing.T) {
	// 获取内容列表
	type ContentGetListInput struct {
		Type       string
		CategoryId uint `p:"cate"`
		Page       int
		Size       int
		Sort       int
		UserId     uint
	}
	// 获取内容列表
	type ContentGetListReq struct {
		ContentGetListInput
		CategoryId uint `p:"cate"`
		Page       int  `d:"1"  v:"min:0#分页号码错误"`
		Size       int  `d:"10" v:"max:50#分页数量最大50条"`
	}

	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var req *ContentGetListReq
		if err := r.Parse(&req); err != nil {
			r.Response.Write(err)
		} else {
			r.Response.Write(req.ContentGetListInput)
		}
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/parse?cate=1&page=2&size=10"), `{"Type":"","CategoryId":1,"Page":2,"Size":10,"Sort":0,"UserId":0}`)
	})
}
