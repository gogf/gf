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

type NamesObject struct{}

func (o *NamesObject) ShowName(r *ghttp.Request) {
	r.Response.Write("Object Show Name")
}

func Test_NameToUri_FullName(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.URI_TYPE_FULLNAME)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/NamesObject"), "Not Found")
		t.Assert(client.GetContent("/NamesObject/ShowName"), "Object Show Name")
	})
}

func Test_NameToUri_AllLower(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/NamesObject"), "Not Found")
		t.Assert(client.GetContent("/namesobject/showname"), "Object Show Name")
	})
}

func Test_NameToUri_Camel(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.URI_TYPE_CAMEL)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/NamesObject"), "Not Found")
		t.Assert(client.GetContent("/namesObject/showName"), "Object Show Name")
	})
}

func Test_NameToUri_Default(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.URI_TYPE_DEFAULT)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/NamesObject"), "Not Found")
		t.Assert(client.GetContent("/names-object/show-name"), "Object Show Name")
	})
}
