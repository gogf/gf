// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

type NamesObject struct{}

func (o *NamesObject) ShowName(r *ghttp.Request) {
	r.Response.Write("Object Show Name")
}

func Test_NameToUri_FullName(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_FULLNAME)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/NamesObject"), "Not Found")
		gtest.Assert(client.GetContent("/NamesObject/ShowName"), "Object Show Name")
	})
}

func Test_NameToUri_AllLower(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_ALLLOWER)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/NamesObject"), "Not Found")
		gtest.Assert(client.GetContent("/namesobject/showname"), "Object Show Name")
	})
}

func Test_NameToUri_Camel(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_CAMEL)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/NamesObject"), "Not Found")
		gtest.Assert(client.GetContent("/namesObject/showName"), "Object Show Name")
	})
}

func Test_NameToUri_Default(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.SetNameToUriType(ghttp.NAME_TO_URI_TYPE_DEFAULT)
	s.BindObject("/{.struct}/{.method}", new(NamesObject))
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/NamesObject"), "Not Found")
		gtest.Assert(client.GetContent("/names-object/show-name"), "Object Show Name")
	})
}
