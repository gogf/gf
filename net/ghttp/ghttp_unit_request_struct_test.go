// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/util/gvalid"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Params_Parse(t *testing.T) {
	type User struct {
		Id   int
		Name string
		Map  map[string]interface{}
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var user *User
		if err := r.Parse(&user); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(user.Map["id"], user.Map["score"])
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.PostContent(ctx, "/parse", `{"id":1,"name":"john","map":{"id":1,"score":100}}`), `1100`)
	})
}

func Test_Params_ParseQuery(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse-query", func(r *ghttp.Request) {
		var user *User
		if err := r.ParseQuery(&user); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(user.Id, user.Name)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.GetContent(ctx, "/parse-query"), `0`)
		t.Assert(c.GetContent(ctx, "/parse-query?id=1&name=john"), `1john`)
		t.Assert(c.PostContent(ctx, "/parse-query"), `0`)
		t.Assert(c.PostContent(ctx, "/parse-query", g.Map{
			"id":   1,
			"name": "john",
		}), `0`)
	})
}

func Test_Params_ParseForm(t *testing.T) {
	type User struct {
		Id   int
		Name string
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse-form", func(r *ghttp.Request) {
		var user *User
		if err := r.ParseForm(&user); err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(user.Id, user.Name)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.GetContent(ctx, "/parse-form"), `0`)
		t.Assert(c.GetContent(ctx, "/parse-form", g.Map{
			"id":   1,
			"name": "john",
		}), 0)
		t.Assert(c.PostContent(ctx, "/parse-form"), `0`)
		t.Assert(c.PostContent(ctx, "/parse-form", g.Map{
			"id":   1,
			"name": "john",
		}), `1john`)
	})
}

func Test_Params_ComplexJsonStruct(t *testing.T) {
	type ItemEnv struct {
		Type  string
		Key   string
		Value string
		Brief string
	}

	type ItemProbe struct {
		Type           string
		Port           int
		Path           string
		Brief          string
		Period         int
		InitialDelay   int
		TimeoutSeconds int
	}

	type ItemKV struct {
		Key   string
		Value string
	}

	type ItemPort struct {
		Port  int
		Type  string
		Alias string
		Brief string
	}

	type ItemMount struct {
		Type    string
		DstPath string
		Src     string
		SrcPath string
		Brief   string
	}

	type SaveRequest struct {
		AppId          uint
		Name           string
		Type           string
		Cluster        string
		Replicas       uint
		ContainerName  string
		ContainerImage string
		VersionTag     string
		Namespace      string
		Id             uint
		Status         uint
		Metrics        string
		InitImage      string
		CpuRequest     uint
		CpuLimit       uint
		MemRequest     uint
		MemLimit       uint
		MeshEnabled    uint
		ContainerPorts []ItemPort
		Labels         []ItemKV
		NodeSelector   []ItemKV
		EnvReserve     []ItemKV
		EnvGlobal      []ItemEnv
		EnvContainer   []ItemEnv
		Mounts         []ItemMount
		LivenessProbe  ItemProbe
		ReadinessProbe ItemProbe
	}

	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			var data *SaveRequest
			if err := r.Parse(&data); err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(data)
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
		content := `
{
    "app_id": 5,
    "cluster": "test",
    "container_image": "nginx",
    "container_name": "test",
    "container_ports": [
        {
            "alias": "别名",
            "brief": "描述",
            "port": 80,
            "type": "tcp"
        }
    ],
    "cpu_limit": 100,
    "cpu_request": 10,
    "create_at": "2020-10-10 12:00:00",
    "creator": 1,
    "env_container": [
        {
            "brief": "用户环境变量",
            "key": "NAME",
            "type": "string",
            "value": "john"
        }
    ],
    "env_global": [
        {
            "brief": "数据数量",
            "key": "NUMBER",
            "type": "string",
            "value": "1"
        }
    ],
    "env_reserve": [
        {
            "key": "NODE_IP",
            "value": "status.hostIP"
        }
    ],
    "liveness_probe": {
        "brief": "存活探针",
        "initial_delay": 10,
        "path": "",
        "period": 5,
        "port": 80,
        "type": "tcpSocket"
    },
    "readiness_probe": {
        "brief": "就绪探针",
        "initial_delay": 10,
        "path": "",
        "period": 5,
        "port": 80,
        "type": "tcpSocket"
    },
    "id": 0,
    "init_image": "",
    "labels": [
        {
            "key": "app",
            "value": "test"
        }
    ],
    "mem_limit": 1000,
    "mem_request": 100,
    "mesh_enabled": 0,
    "metrics": "",
    "mounts": [],
    "name": "test",
    "namespace": "test",
    "node_selector": [
        {
            "key": "group",
            "value": "app"
        }
    ],
    "replicas": 1,
    "type": "test",
    "update_at": "2020-10-10 12:00:00",
    "version_tag": "test"
}
`
		t.Assert(client.PostContent(ctx, "/parse", content), `{"AppId":5,"Name":"test","Type":"test","Cluster":"test","Replicas":1,"ContainerName":"test","ContainerImage":"nginx","VersionTag":"test","Namespace":"test","Id":0,"Status":0,"Metrics":"","InitImage":"","CpuRequest":10,"CpuLimit":100,"MemRequest":100,"MemLimit":1000,"MeshEnabled":0,"ContainerPorts":[{"Port":80,"Type":"tcp","Alias":"别名","Brief":"描述"}],"Labels":[{"Key":"app","Value":"test"}],"NodeSelector":[{"Key":"group","Value":"app"}],"EnvReserve":[{"Key":"NODE_IP","Value":"status.hostIP"}],"EnvGlobal":[{"Type":"string","Key":"NUMBER","Value":"1","Brief":"数据数量"}],"EnvContainer":[{"Type":"string","Key":"NAME","Value":"john","Brief":"用户环境变量"}],"Mounts":[],"LivenessProbe":{"Type":"tcpSocket","Port":80,"Path":"","Brief":"存活探针","Period":5,"InitialDelay":10,"TimeoutSeconds":0},"ReadinessProbe":{"Type":"tcpSocket","Port":80,"Path":"","Brief":"就绪探针","Period":5,"InitialDelay":10,"TimeoutSeconds":0}}`)
	})
}

func Test_Params_Parse_Attr_Pointer1(t *testing.T) {
	type User struct {
		Id   *int
		Name *string
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse1", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			var user *User
			if err := r.Parse(&user); err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(user.Id, user.Name)
		}
	})
	s.BindHandler("/parse2", func(r *ghttp.Request) {
		if m := r.GetMap(); len(m) > 0 {
			var user = new(User)
			if err := r.Parse(user); err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(user.Id, user.Name)
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
		t.Assert(client.PostContent(ctx, "/parse1", `{"id":1,"name":"john"}`), `1john`)
		t.Assert(client.PostContent(ctx, "/parse2", `{"id":1,"name":"john"}`), `1john`)
		t.Assert(client.PostContent(ctx, "/parse2?id=1&name=john"), `1john`)
		t.Assert(client.PostContent(ctx, "/parse2", `id=1&name=john`), `1john`)
	})
}

func Test_Params_Parse_Attr_Pointer2(t *testing.T) {
	type User struct {
		Id *int `v:"required"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/parse", func(r *ghttp.Request) {
		var user *User
		if err := r.Parse(&user); err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(user.Id)
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.PostContent(ctx, "/parse"), `The Id field is required`)
		t.Assert(client.PostContent(ctx, "/parse?id=1"), `1`)
	})
}

// It does not support this kind of converting yet.
//func Test_Params_Parse_Attr_SliceSlice(t *testing.T) {
//	type User struct {
//		Id     int
//		Name   string
//		Scores [][]int
//	}
//	p, _ := ports.PopRand()
//	s := g.Server(p)
//	s.BindHandler("/parse", func(r *ghttp.Request) {
//		if m := r.GetMap(); len(m) > 0 {
//			var user *User
//			if err := r.Parse(&user); err != nil {
//				r.Response.WriteExit(err)
//			}
//			r.Response.WriteExit(user.Scores)
//		}
//	})
//	s.SetPort(p)
//	s.SetDumpRouterMap(false)
//	s.Start()
//	defer s.Shutdown()
//
//	time.Sleep(100 * time.Millisecond)
//	gtest.C(t, func(t *gtest.T) {
//		client := g.Client()
//		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
//		t.Assert(client.PostContent(ctx, "/parse", `{"id":1,"name":"john","scores":[[1,2,3]]}`), `1100`)
//	})
//}

func Test_Params_Struct(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Time  *time.Time
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"password2 @required|length:2,20|password3#||密码强度不足"`
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
			if err := gvalid.CheckStruct(r.Context(), user, nil); err != nil {
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
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent(ctx, "/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent(ctx, "/struct1", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent(ctx, "/struct2", `id=1&name=john&password1=123&password2=456`), `1john123456`)
		t.Assert(client.PostContent(ctx, "/struct2", ``), ``)
		t.Assert(client.PostContent(ctx, "/struct-valid", `id=1&name=john&password1=123&password2=0`), "The password2 value `0` length must be between 2 and 20; 密码强度不足")
		t.Assert(client.PostContent(ctx, "/parse", `id=1&name=john&password1=123&password2=0`), "The password2 value `0` length must be between 2 and 20; 密码强度不足")
		t.Assert(client.PostContent(ctx, "/parse", `{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}`), `1john123Abc!@#123Abc!@#`)
	})
}

func Test_Params_Structs(t *testing.T) {
	type User struct {
		Id    int
		Name  string
		Time  *time.Time
		Pass1 string `p:"password1"`
		Pass2 string `p:"password2" v:"password2 @required|length:2,20|password3#||密码强度不足"`
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
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.PostContent(ctx,
			"/parse1",
			`[{"id":1,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}, {"id":2,"name":"john","password1":"123Abc!@#","password2":"123Abc!@#"}]`),
			`12`,
		)
	})
}

func Test_Params_Struct_Validation(t *testing.T) {
	type User struct {
		Id   int    `v:"required"`
		Name string `v:"name@required-with:id"`
	}
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			var (
				err  error
				user *User
			)
			err = r.Parse(&user)
			if err != nil {
				r.Response.WriteExit(err)
			}
			r.Response.WriteExit(user.Id, user.Name)
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(c.GetContent(ctx, "/", ``), `The Id field is required`)
		t.Assert(c.GetContent(ctx, "/", `id=1&name=john`), `1john`)
		t.Assert(c.PostContent(ctx, "/", `id=1&name=john&password1=123&password2=456`), `1john`)
		t.Assert(c.PostContent(ctx, "/", `id=1`), `The name field is required`)
	})
}
