// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

// 执行对象
type GroupObject struct{}

func (o *GroupObject) Init(r *ghttp.Request) {
	r.Response.Write("1")
}

func (o *GroupObject) Shut(r *ghttp.Request) {
	r.Response.Write("2")
}

func (o *GroupObject) Index(r *ghttp.Request) {
	r.Response.Write("Object Index")
}

func (o *GroupObject) Show(r *ghttp.Request) {
	r.Response.Write("Object Show")
}

func (o *GroupObject) Delete(r *ghttp.Request) {
	r.Response.Write("Object Delete")
}

func Handler(r *ghttp.Request) {
	r.Response.Write("Handler")
}

func Test_Router_GroupBasic1(t *testing.T) {
	s := g.Server(guid.S())
	obj := new(GroupObject)
	// 分组路由方法注册
	group := s.Group("/api")
	group.ALL("/handler", Handler)
	group.ALL("/obj", obj)
	group.GET("/obj/my-show", obj, "Show")
	group.REST("/obj/rest", obj)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/api/handler"), "Handler")

		t.Assert(client.GetContent(ctx, "/api/obj"), "1Object Index2")
		t.Assert(client.GetContent(ctx, "/api/obj/"), "1Object Index2")
		t.Assert(client.GetContent(ctx, "/api/obj/index"), "1Object Index2")
		t.Assert(client.GetContent(ctx, "/api/obj/delete"), "1Object Delete2")
		t.Assert(client.GetContent(ctx, "/api/obj/my-show"), "1Object Show2")
		t.Assert(client.GetContent(ctx, "/api/obj/show"), "1Object Show2")
		t.Assert(client.DeleteContent(ctx, "/api/obj/rest"), "1Object Delete2")

		t.Assert(client.DeleteContent(ctx, "/ThisDoesNotExist"), "Not Found")
		t.Assert(client.DeleteContent(ctx, "/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_GroupBuildInVar(t *testing.T) {
	s := g.Server(guid.S())
	obj := new(GroupObject)
	// 分组路由方法注册
	group := s.Group("/api")
	group.ALL("/{.struct}/{.method}", obj)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/api/group-object/index"), "1Object Index2")
		t.Assert(client.GetContent(ctx, "/api/group-object/delete"), "1Object Delete2")
		t.Assert(client.GetContent(ctx, "/api/group-object/show"), "1Object Show2")

		t.Assert(client.DeleteContent(ctx, "/ThisDoesNotExist"), "Not Found")
		t.Assert(client.DeleteContent(ctx, "/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_Group_Methods(t *testing.T) {
	s := g.Server(guid.S())
	obj := new(GroupObject)
	group := s.Group("/")
	group.ALL("/obj", obj, "Show, Delete")
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(client.GetContent(ctx, "/obj/show"), "1Object Show2")
		t.Assert(client.GetContent(ctx, "/obj/delete"), "1Object Delete2")
	})
}

func Test_Router_Group_MultiServer(t *testing.T) {
	s1 := g.Server(guid.S())
	s2 := g.Server(guid.S())
	s1.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/post", func(r *ghttp.Request) {
			r.Response.Write("post1")
		})
	})
	s2.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/post", func(r *ghttp.Request) {
			r.Response.Write("post2")
		})
	})
	s1.SetDumpRouterMap(false)
	s2.SetDumpRouterMap(false)
	gtest.Assert(s1.Start(), nil)
	gtest.Assert(s2.Start(), nil)
	defer s1.Shutdown()
	defer s2.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c1 := g.Client()
		c1.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s1.GetListenedPort()))
		c2 := g.Client()
		c2.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s2.GetListenedPort()))
		t.Assert(c1.PostContent(ctx, "/post"), "post1")
		t.Assert(c2.PostContent(ctx, "/post"), "post2")
	})
}

func Test_Router_Group_Map(t *testing.T) {
	testFuncGet := func(r *ghttp.Request) {
		r.Response.Write("get")
	}
	testFuncPost := func(r *ghttp.Request) {
		r.Response.Write("post")
	}
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Map(map[string]any{
			"Get: /test": testFuncGet,
			"Post:/test": testFuncPost,
		})
	})
	//s.SetDumpRouterMap(false)
	gtest.Assert(s.Start(), nil)
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/test"), "get")
		t.Assert(c.PostContent(ctx, "/test"), "post")
	})
}

type SafeBuffer struct {
	buffer *bytes.Buffer
	mu     sync.Mutex
}

func (b *SafeBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.Write(p)
}

func (b *SafeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.String()
}

func Test_Router_OverWritten(t *testing.T) {
	var (
		s   = g.Server(guid.S())
		obj = new(GroupObject)
		buf = &SafeBuffer{
			buffer: bytes.NewBuffer(nil),
			mu:     sync.Mutex{},
		}
		logger = glog.NewWithWriter(buf)
	)
	logger.SetStdoutPrint(false)
	s.SetLogger(logger)
	s.SetRouteOverWrite(true)
	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.ALLMap(g.Map{
			"/obj": obj,
		})
		group.ALLMap(g.Map{
			"/obj": obj,
		})
	})
	s.Start()
	defer s.Shutdown()

	dumpContent := buf.String()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Count(dumpContent, `/api/obj `), 1)
		t.Assert(gstr.Count(dumpContent, `/api/obj/index`), 1)
		t.Assert(gstr.Count(dumpContent, `/api/obj/show`), 1)
		t.Assert(gstr.Count(dumpContent, `/api/obj/delete`), 1)
	})
}
