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

type DomainObject struct{}

func (o *DomainObject) Init(r *ghttp.Request) {
	r.Response.Write("1")
}

func (o *DomainObject) Shut(r *ghttp.Request) {
	r.Response.Write("2")
}

func (o *DomainObject) Index(r *ghttp.Request) {
	r.Response.Write("Object Index")
}

func (o *DomainObject) Show(r *ghttp.Request) {
	r.Response.Write("Object Show")
}

func (o *DomainObject) Info(r *ghttp.Request) {
	r.Response.Write("Object Info")
}

func Test_Router_DomainObject1(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindObject("/", new(DomainObject))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/init"), "Not Found")
		t.Assert(client.GetContent("/shut"), "Not Found")
		t.Assert(client.GetContent("/index"), "Not Found")
		t.Assert(client.GetContent("/show"), "Not Found")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "1Object Index2")
		t.Assert(client.GetContent("/init"), "Not Found")
		t.Assert(client.GetContent("/shut"), "Not Found")
		t.Assert(client.GetContent("/index"), "1Object Index2")
		t.Assert(client.GetContent("/show"), "1Object Show2")
		t.Assert(client.GetContent("/info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})

	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "1Object Index2")
		t.Assert(client.GetContent("/init"), "Not Found")
		t.Assert(client.GetContent("/shut"), "Not Found")
		t.Assert(client.GetContent("/index"), "1Object Index2")
		t.Assert(client.GetContent("/show"), "1Object Show2")
		t.Assert(client.GetContent("/info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}

func Test_Router_DomainObject2(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindObject("/object", new(DomainObject), "Show, Info")
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "Not Found")
		t.Assert(client.GetContent("/object/info"), "Not Found")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "1Object Show2")
		t.Assert(client.GetContent("/object/info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "1Object Show2")
		t.Assert(client.GetContent("/object/info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}

func Test_Router_DomainObjectMethod(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Domain("localhost, local").BindObjectMethod("/object-info", new(DomainObject), "Info")
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "Not Found")
		t.Assert(client.GetContent("/object/info"), "Not Found")
		t.Assert(client.GetContent("/object-info"), "Not Found")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "Not Found")
		t.Assert(client.GetContent("/object/info"), "Not Found")
		t.Assert(client.GetContent("/object-info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/object"), "Not Found")
		t.Assert(client.GetContent("/object/init"), "Not Found")
		t.Assert(client.GetContent("/object/shut"), "Not Found")
		t.Assert(client.GetContent("/object/index"), "Not Found")
		t.Assert(client.GetContent("/object/show"), "Not Found")
		t.Assert(client.GetContent("/object/info"), "Not Found")
		t.Assert(client.GetContent("/object-info"), "1Object Info2")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
