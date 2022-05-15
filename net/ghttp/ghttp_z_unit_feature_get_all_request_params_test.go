package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
	"time"
)

func Test_GetAll_Request_Params(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/get_all_params", func(r *ghttp.Request) {
		r.Response.Write(r.GetAllRequestParams())
	})
	s.BindHandler("/get_all_params/:id", func(r *ghttp.Request) {
		r.Response.Write(r.GetAllRequestParams())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		// GET
		t.Assert(client.GetContent(ctx, "/get_all_params", "array[]=1&array[]=2"), `{"array":["1","2"]}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "slice=1&slice=2"), `{"slice":"2"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "bool=1"), `{"bool":"1"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "bool=0"), `{"bool":"0"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "float32=0.11"), `{"float32":"0.11"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "float64=0.22"), `{"float64":"0.22"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "int=-10000"), `{"int":"-10000"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "int=10000"), `{"int":"10000"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "uint=10000"), `{"uint":"10000"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "uint=9"), `{"uint":"9"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "string=key"), `{"string":"key"}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "map[a]=1&map[b]=2"), `{"map":{"a":"1","b":"2"}}`)
		t.Assert(client.GetContent(ctx, "/get_all_params", "a=1&b=2"), `{"a":"1","b":"2"}`)
		//
		// PUT
		t.Assert(client.PutContent(ctx, "/get_all_params", "array[]=1&array[]=2"), `{"array":["1","2"]}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "slice=1&slice=2"), `{"slice":"2"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "bool=1"), `{"bool":"1"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "bool=0"), `{"bool":"0"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "float32=0.11"), `{"float32":"0.11"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "float64=0.22"), `{"float64":"0.22"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "int=-10000"), `{"int":"-10000"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "int=10000"), `{"int":"10000"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "uint=10000"), `{"uint":"10000"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "uint=9"), `{"uint":"9"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "string=key"), `{"string":"key"}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "map[a]=1&map[b]=2"), `{"map":{"a":"1","b":"2"}}`)
		t.Assert(client.PutContent(ctx, "/get_all_params", "a=1&b=2"), `{"a":"1","b":"2"}`)
		//
		// DELETE
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "array[]=1&array[]=2"), `{"array":["1","2"]}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "slice=1&slice=2"), `{"slice":"2"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "bool=1"), `{"bool":"1"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "bool=0"), `{"bool":"0"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "float32=0.11"), `{"float32":"0.11"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "float64=0.22"), `{"float64":"0.22"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "int=-10000"), `{"int":"-10000"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "int=10000"), `{"int":"10000"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "uint=10000"), `{"uint":"10000"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "uint=9"), `{"uint":"9"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "string=key"), `{"string":"key"}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "map[a]=1&map[b]=2"), `{"map":{"a":"1","b":"2"}}`)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", "a=1&b=2"), `{"a":"1","b":"2"}`)
		//
		// Form
		t.Assert(client.PostContent(ctx, "/get_all_params", "array[]=1&array[]=2"), `{"array":["1","2"]}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "slice=1&slice=2"), `{"slice":"2"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "bool=1"), `{"bool":"1"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "bool=0"), `{"bool":"0"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "float32=0.11"), `{"float32":"0.11"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "float64=0.22"), `{"float64":"0.22"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "int=-10000"), `{"int":"-10000"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "int=10000"), `{"int":"10000"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "uint=10000"), `{"uint":"10000"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "uint=9"), `{"uint":"9"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "string=key"), `{"string":"key"}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "map[a]=1&map[b]=2"), `{"map":{"a":"1","b":"2"}}`)
		t.Assert(client.PostContent(ctx, "/get_all_params", "a=1&b=2"), `{"a":"1","b":"2"}`)
		//
		// Json
		t.Assert(client.PostContent(ctx, "/get_all_params", `{"id":1,"name":"john"}`), `{"id":1,"name":"john"}`)
		//
		// Empty
		t.Assert(client.GetContent(ctx, "/get_all_params", ``), nil)
		t.Assert(client.PutContent(ctx, "/get_all_params", ``), nil)
		t.Assert(client.DeleteContent(ctx, "/get_all_params", ``), nil)
		t.Assert(client.PostContent(ctx, "/get_all_params", ``), nil)
		//
		// Router
		t.Assert(client.GetContent(ctx, "/get_all_params/1", ``), `{"id":"1"}`)
	})
}
