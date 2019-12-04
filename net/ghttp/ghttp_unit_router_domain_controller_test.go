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

type DomainController struct {
	gmvc.Controller
}

func (c *DomainController) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.Response.Write("1")
}

func (c *DomainController) Shut() {
	c.Response.Write("2")
}

func (c *DomainController) Index() {
	c.Response.Write("Controller Index")
}

func (c *DomainController) Show() {
	c.Response.Write("Controller Show")
}

func (c *DomainController) Info() {
	c.Response.Write("Controller Info")
}

func Test_Router_DomainController1(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindController("/", new(DomainController))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/init"), "Not Found")
		gtest.Assert(client.GetContent("/shut"), "Not Found")
		gtest.Assert(client.GetContent("/index"), "Not Found")
		gtest.Assert(client.GetContent("/show"), "Not Found")
		gtest.Assert(client.GetContent("/info"), "Not Found")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "1Controller Index2")
		gtest.Assert(client.GetContent("/init"), "Not Found")
		gtest.Assert(client.GetContent("/shut"), "Not Found")
		gtest.Assert(client.GetContent("/index"), "1Controller Index2")
		gtest.Assert(client.GetContent("/show"), "1Controller Show2")
		gtest.Assert(client.GetContent("/info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "1Controller Index2")
		gtest.Assert(client.GetContent("/init"), "Not Found")
		gtest.Assert(client.GetContent("/shut"), "Not Found")
		gtest.Assert(client.GetContent("/index"), "1Controller Index2")
		gtest.Assert(client.GetContent("/show"), "1Controller Show2")
		gtest.Assert(client.GetContent("/info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}

func Test_Router_DomainController2(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindController("/controller", new(DomainController), "Show, Info")
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "Not Found")
		gtest.Assert(client.GetContent("/controller/info"), "Not Found")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "1Controller Show2")
		gtest.Assert(client.GetContent("/controller/info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "1Controller Show2")
		gtest.Assert(client.GetContent("/controller/info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}

func Test_Router_DomainControllerMethod(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindControllerMethod("/controller-info", new(DomainController), "Info")
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "Not Found")
		gtest.Assert(client.GetContent("/controller/info"), "Not Found")
		gtest.Assert(client.GetContent("/controller-info"), "Not Found")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "Not Found")
		gtest.Assert(client.GetContent("/controller/info"), "Not Found")
		gtest.Assert(client.GetContent("/controller-info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/controller"), "Not Found")
		gtest.Assert(client.GetContent("/controller/init"), "Not Found")
		gtest.Assert(client.GetContent("/controller/shut"), "Not Found")
		gtest.Assert(client.GetContent("/controller/index"), "Not Found")
		gtest.Assert(client.GetContent("/controller/show"), "Not Found")
		gtest.Assert(client.GetContent("/controller/info"), "Not Found")
		gtest.Assert(client.GetContent("/controller-info"), "1Controller Info2")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
