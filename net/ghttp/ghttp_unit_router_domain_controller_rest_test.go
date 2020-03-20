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

type DomainControllerRest struct {
	gmvc.Controller
}

func (c *DomainControllerRest) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.Response.Write("1")
}

func (c *DomainControllerRest) Shut() {
	c.Response.Write("2")
}

func (c *DomainControllerRest) Get() {
	c.Response.Write("Controller Get")
}

func (c *DomainControllerRest) Put() {
	c.Response.Write("Controller Put")
}

func (c *DomainControllerRest) Post() {
	c.Response.Write("Controller Post")
}

func (c *DomainControllerRest) Delete() {
	c.Response.Write("Controller Delete")
}

func (c *DomainControllerRest) Patch() {
	c.Response.Write("Controller Patch")
}

func (c *DomainControllerRest) Options() {
	c.Response.Write("Controller Options")
}

func (c *DomainControllerRest) Head() {
	c.Response.Header().Set("head-ok", "1")
}

// 控制器注册测试
func Test_Router_DomainControllerRest(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindControllerRest("/", new(DomainControllerRest))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.PutContent("/"), "Not Found")
		t.Assert(client.PostContent("/"), "Not Found")
		t.Assert(client.DeleteContent("/"), "Not Found")
		t.Assert(client.PatchContent("/"), "Not Found")
		t.Assert(client.OptionsContent("/"), "Not Found")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "1Controller Get2")
		t.Assert(client.PutContent("/"), "1Controller Put2")
		t.Assert(client.PostContent("/"), "1Controller Post2")
		t.Assert(client.DeleteContent("/"), "1Controller Delete2")
		t.Assert(client.PatchContent("/"), "1Controller Patch2")
		t.Assert(client.OptionsContent("/"), "1Controller Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "1")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "1Controller Get2")
		t.Assert(client.PutContent("/"), "1Controller Put2")
		t.Assert(client.PostContent("/"), "1Controller Post2")
		t.Assert(client.DeleteContent("/"), "1Controller Delete2")
		t.Assert(client.PatchContent("/"), "1Controller Patch2")
		t.Assert(client.OptionsContent("/"), "1Controller Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "1")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
