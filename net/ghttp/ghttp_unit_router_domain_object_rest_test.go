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

type DomainObjectRest struct{}

func (o *DomainObjectRest) Init(r *ghttp.Request) {
	r.Response.Write("1")
}

func (o *DomainObjectRest) Shut(r *ghttp.Request) {
	r.Response.Write("2")
}

func (o *DomainObjectRest) Get(r *ghttp.Request) {
	r.Response.Write("Object Get")
}

func (o *DomainObjectRest) Put(r *ghttp.Request) {
	r.Response.Write("Object Put")
}

func (o *DomainObjectRest) Post(r *ghttp.Request) {
	r.Response.Write("Object Post")
}

func (o *DomainObjectRest) Delete(r *ghttp.Request) {
	r.Response.Write("Object Delete")
}

func (o *DomainObjectRest) Patch(r *ghttp.Request) {
	r.Response.Write("Object Patch")
}

func (o *DomainObjectRest) Options(r *ghttp.Request) {
	r.Response.Write("Object Options")
}

func (o *DomainObjectRest) Head(r *ghttp.Request) {
	r.Response.Header().Set("head-ok", "1")
}

func Test_Router_DomainObjectRest(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindObjectRest("/", new(DomainObjectRest))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.PutContent("/"), "Not Found")
		gtest.Assert(client.PostContent("/"), "Not Found")
		gtest.Assert(client.DeleteContent("/"), "Not Found")
		gtest.Assert(client.PatchContent("/"), "Not Found")
		gtest.Assert(client.OptionsContent("/"), "Not Found")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		gtest.Assert(err, nil)
		gtest.Assert(resp1.Header.Get("head-ok"), "")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "1Object Get2")
		gtest.Assert(client.PutContent("/"), "1Object Put2")
		gtest.Assert(client.PostContent("/"), "1Object Post2")
		gtest.Assert(client.DeleteContent("/"), "1Object Delete2")
		gtest.Assert(client.PatchContent("/"), "1Object Patch2")
		gtest.Assert(client.OptionsContent("/"), "1Object Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		gtest.Assert(err, nil)
		gtest.Assert(resp1.Header.Get("head-ok"), "1")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "1Object Get2")
		gtest.Assert(client.PutContent("/"), "1Object Put2")
		gtest.Assert(client.PostContent("/"), "1Object Post2")
		gtest.Assert(client.DeleteContent("/"), "1Object Delete2")
		gtest.Assert(client.PatchContent("/"), "1Object Patch2")
		gtest.Assert(client.OptionsContent("/"), "1Object Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		gtest.Assert(err, nil)
		gtest.Assert(resp1.Header.Get("head-ok"), "1")
		gtest.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
