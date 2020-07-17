// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_Params_Page(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/list", func(r *ghttp.Request) {
			page := r.GetPage(5, 2)
			r.Response.Write(page.GetContent(4))
		})
		group.GET("/list/{page}.html", func(r *ghttp.Request) {
			page := r.GetPage(5, 2)
			r.Response.Write(page.GetContent(4))
		})
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/list"), `<span class="GPageSpan">首页</span><span class="GPageSpan">上一页</span><span class="GPageSpan">1</span><a class="GPageLink" href="/list?page=2" title="2">2</a><a class="GPageLink" href="/list?page=3" title="3">3</a><a class="GPageLink" href="/list?page=2" title="">下一页</a><a class="GPageLink" href="/list?page=3" title="">尾页</a>`)
		t.Assert(client.GetContent("/list?page=3"), `<a class="GPageLink" href="/list?page=1" title="">首页</a><a class="GPageLink" href="/list?page=2" title="">上一页</a><a class="GPageLink" href="/list?page=1" title="1">1</a><a class="GPageLink" href="/list?page=2" title="2">2</a><span class="GPageSpan">3</span><span class="GPageSpan">下一页</span><span class="GPageSpan">尾页</span>`)

		t.Assert(client.GetContent("/list/1.html"), `<span class="GPageSpan">首页</span><span class="GPageSpan">上一页</span><span class="GPageSpan">1</span><a class="GPageLink" href="/list/2.html" title="2">2</a><a class="GPageLink" href="/list/3.html" title="3">3</a><a class="GPageLink" href="/list/2.html" title="">下一页</a><a class="GPageLink" href="/list/3.html" title="">尾页</a>`)
		t.Assert(client.GetContent("/list/3.html"), `<a class="GPageLink" href="/list/1.html" title="">首页</a><a class="GPageLink" href="/list/2.html" title="">上一页</a><a class="GPageLink" href="/list/1.html" title="1">1</a><a class="GPageLink" href="/list/2.html" title="2">2</a><span class="GPageSpan">3</span><span class="GPageSpan">下一页</span><span class="GPageSpan">尾页</span>`)
	})
}
