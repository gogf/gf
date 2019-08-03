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

func Test_Router_MiddleWare_Basic1(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/test/test", func(r *ghttp.Request) {
		r.Response.Write("test")
	})
	s.AddMiddleWare(func(r *ghttp.Request) {
		r.Response.Write("1")
		r.MiddleWare.Next()
		r.Response.Write("2")
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "12")
		gtest.Assert(client.GetContent("/test/test"), "1test2")
	})
}
