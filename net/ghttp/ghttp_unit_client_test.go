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

func Test_Client_Basic(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		url := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := ghttp.NewClient()
		client.SetPrefix(url)

		gtest.Assert(ghttp.GetContent(""), ``)
		gtest.Assert(client.GetContent("/hello"), `hello`)

		_, err := ghttp.Post("")
		gtest.AssertNE(err, nil)
	})
}
