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
	p, _ := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindObjectRest("/", new(DomainObjectRest))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
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
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		t.Assert(client.GetContent("/"), "1Object Get2")
		t.Assert(client.PutContent("/"), "1Object Put2")
		t.Assert(client.PostContent("/"), "1Object Post2")
		t.Assert(client.DeleteContent("/"), "1Object Delete2")
		t.Assert(client.PatchContent("/"), "1Object Patch2")
		t.Assert(client.OptionsContent("/"), "1Object Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "1")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		t.Assert(client.GetContent("/"), "1Object Get2")
		t.Assert(client.PutContent("/"), "1Object Put2")
		t.Assert(client.PostContent("/"), "1Object Post2")
		t.Assert(client.DeleteContent("/"), "1Object Delete2")
		t.Assert(client.PatchContent("/"), "1Object Patch2")
		t.Assert(client.OptionsContent("/"), "1Object Options2")
		resp1, err := client.Head("/")
		if err == nil {
			defer resp1.Close()
		}
		t.Assert(err, nil)
		t.Assert(resp1.Header.Get("head-ok"), "1")
		t.Assert(client.GetContent("/none-exist"), "Not Found")
	})
}
