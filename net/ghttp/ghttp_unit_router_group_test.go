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
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
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

// 控制器
type GroupController struct {
	gmvc.Controller
}

func (c *GroupController) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.Response.Write("1")
}

func (c *GroupController) Shut() {
	c.Response.Write("2")
}

func (c *GroupController) Index() {
	c.Response.Write("Controller Index")
}

func (c *GroupController) Show() {
	c.Response.Write("Controller Show")
}

func (c *GroupController) Post() {
	c.Response.Write("Controller Post")
}

func Handler(r *ghttp.Request) {
	r.Response.Write("Handler")
}

func Test_Router_GroupBasic1(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	obj := new(GroupObject)
	ctl := new(GroupController)
	// 分组路由方法注册
	group := s.Group("/api")
	group.ALL("/handler", Handler)
	group.ALL("/ctl", ctl)
	group.GET("/ctl/my-show", ctl, "Show")
	group.REST("/ctl/rest", ctl)
	group.ALL("/obj", obj)
	group.GET("/obj/my-show", obj, "Show")
	group.REST("/obj/rest", obj)
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/api/handler"), "Handler")

		t.Assert(client.GetContent("/api/ctl"), "1Controller Index2")
		t.Assert(client.GetContent("/api/ctl/"), "1Controller Index2")
		t.Assert(client.GetContent("/api/ctl/index"), "1Controller Index2")
		t.Assert(client.GetContent("/api/ctl/my-show"), "1Controller Show2")
		t.Assert(client.GetContent("/api/ctl/post"), "1Controller Post2")
		t.Assert(client.GetContent("/api/ctl/show"), "1Controller Show2")
		t.Assert(client.PostContent("/api/ctl/rest"), "1Controller Post2")

		t.Assert(client.GetContent("/api/obj"), "1Object Index2")
		t.Assert(client.GetContent("/api/obj/"), "1Object Index2")
		t.Assert(client.GetContent("/api/obj/index"), "1Object Index2")
		t.Assert(client.GetContent("/api/obj/delete"), "1Object Delete2")
		t.Assert(client.GetContent("/api/obj/my-show"), "1Object Show2")
		t.Assert(client.GetContent("/api/obj/show"), "1Object Show2")
		t.Assert(client.DeleteContent("/api/obj/rest"), "1Object Delete2")

		t.Assert(client.DeleteContent("/ThisDoesNotExist"), "Not Found")
		t.Assert(client.DeleteContent("/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_GroupBasic2(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	obj := new(GroupObject)
	ctl := new(GroupController)
	// 分组路由批量注册
	s.Group("/api").Bind([]g.Slice{
		{"ALL", "/handler", Handler},
		{"ALL", "/ctl", ctl},
		{"GET", "/ctl/my-show", ctl, "Show"},
		{"REST", "/ctl/rest", ctl},
		{"ALL", "/obj", obj},
		{"GET", "/obj/my-show", obj, "Show"},
		{"REST", "/obj/rest", obj},
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/api/handler"), "Handler")

		t.Assert(client.GetContent("/api/ctl/my-show"), "1Controller Show2")
		t.Assert(client.GetContent("/api/ctl/post"), "1Controller Post2")
		t.Assert(client.GetContent("/api/ctl/show"), "1Controller Show2")
		t.Assert(client.PostContent("/api/ctl/rest"), "1Controller Post2")

		t.Assert(client.GetContent("/api/obj/delete"), "1Object Delete2")
		t.Assert(client.GetContent("/api/obj/my-show"), "1Object Show2")
		t.Assert(client.GetContent("/api/obj/show"), "1Object Show2")
		t.Assert(client.DeleteContent("/api/obj/rest"), "1Object Delete2")

		t.Assert(client.DeleteContent("/ThisDoesNotExist"), "Not Found")
		t.Assert(client.DeleteContent("/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_GroupBuildInVar(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	obj := new(GroupObject)
	ctl := new(GroupController)
	// 分组路由方法注册
	group := s.Group("/api")
	group.ALL("/{.struct}/{.method}", ctl)
	group.ALL("/{.struct}/{.method}", obj)
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/api/group-controller/index"), "1Controller Index2")
		t.Assert(client.GetContent("/api/group-controller/post"), "1Controller Post2")
		t.Assert(client.GetContent("/api/group-controller/show"), "1Controller Show2")

		t.Assert(client.GetContent("/api/group-object/index"), "1Object Index2")
		t.Assert(client.GetContent("/api/group-object/delete"), "1Object Delete2")
		t.Assert(client.GetContent("/api/group-object/show"), "1Object Show2")

		t.Assert(client.DeleteContent("/ThisDoesNotExist"), "Not Found")
		t.Assert(client.DeleteContent("/api/ThisDoesNotExist"), "Not Found")
	})
}

func Test_Router_Group_Mthods(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	obj := new(GroupObject)
	ctl := new(GroupController)
	group := s.Group("/")
	group.ALL("/obj", obj, "Show, Delete")
	group.ALL("/ctl", ctl, "Show, Post")
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/ctl/show"), "1Controller Show2")
		t.Assert(client.GetContent("/ctl/post"), "1Controller Post2")
		t.Assert(client.GetContent("/obj/show"), "1Object Show2")
		t.Assert(client.GetContent("/obj/delete"), "1Object Delete2")
	})
}
