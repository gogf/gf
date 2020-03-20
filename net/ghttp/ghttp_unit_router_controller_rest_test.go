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

type ControllerRest struct {
	gmvc.Controller
}

func (c *ControllerRest) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.Response.Write("1")
}

func (c *ControllerRest) Shut() {
	c.Response.Write("2")
}

func (c *ControllerRest) Get() {
	c.Response.Write("Controller Get")
}

func (c *ControllerRest) Put() {
	c.Response.Write("Controller Put")
}

func (c *ControllerRest) Post() {
	c.Response.Write("Controller Post")
}

func (c *ControllerRest) Delete() {
	c.Response.Write("Controller Delete")
}

func (c *ControllerRest) Head() {
	c.Response.Header().Set("head-ok", "1")
}

// 控制器注册测试
func Test_Router_ControllerRest(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindControllerRest("/", new(ControllerRest))
	s.BindControllerRest("/{.struct}/{.method}", new(ControllerRest))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "1Controller Get2")
		t.Assert(client.PutContent("/"), "1Controller Put2")
		t.Assert(client.PostContent("/"), "1Controller Post2")
		t.Assert(client.DeleteContent("/"), "1Controller Delete2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "1")

		t.Assert(client.GetContent("/controller-rest/get"), "1Controller Get2")
		t.Assert(client.PutContent("/controller-rest/put"), "1Controller Put2")
		t.Assert(client.PostContent("/controller-rest/post"), "1Controller Post2")
		t.Assert(client.DeleteContent("/controller-rest/delete"), "1Controller Delete2")
		resp2, err := client.Head("/controller-rest/head")
		if err == nil {
			defer resp2.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp2.Header.Get("head-ok"), "1")

		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
