// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MPL was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Client_Query(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Write(r.URL.RawQuery)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	var prefix = fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		t.Assert(c.Query(g.Map{
			"page": 1,
			"tags": []string{"go", "web"},
		}).GetContent(ctx, "/query"), "page=1&tags=go&tags=web")
	})
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		// Existing query parameters in the url are kept as is.
		t.Assert(c.Query(g.Map{"page": 2}).GetContent(ctx, "/query?size=10"), "size=10&page=2")
	})
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		// Query parameters are also applied to non-GET requests.
		t.Assert(c.Query(g.Map{"page": 3}).PostContent(ctx, "/query"), "page=3")
	})
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		// The chaining method does not modify the original client.
		_ = c.Query(g.Map{"page": 1})
		t.Assert(c.GetContent(ctx, "/query"), "")
	})
}

func Test_Client_SetQueryMap(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/query", func(r *ghttp.Request) {
		r.Response.Write(r.URL.RawQuery)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	var prefix = fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		c.SetQuery("page", 1)
		c.SetQueryMap(g.Map{"size": 10})
		t.Assert(c.GetContent(ctx, "/query"), "page=1&size=10")
	})
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix)
		// Values are url-encoded by default.
		c.SetQuery("keyword", "go frame")
		t.Assert(c.GetContent(ctx, "/query"), "keyword=go+frame")
	})
	gtest.C(t, func(t *gtest.T) {
		var c = g.Client().SetPrefix(prefix).NoUrlEncode()
		c.SetQuery("keyword", "go,frame")
		t.Assert(c.GetContent(ctx, "/query"), "keyword=go,frame")
	})
}
