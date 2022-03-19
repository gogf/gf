// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	ctx = context.TODO()
)

func init() {
	genv.Set("UNDER_TEST", "1")
}

func Test_GetUrl(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/url", func(r *ghttp.Request) {
		r.Response.Write(r.GetUrl())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/url"), prefix+"/url")
	})
}

func Test_XUrlPath(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test1", func(r *ghttp.Request) {
		r.Response.Write(`test1`)
	})
	s.BindHandler("/test2", func(r *ghttp.Request) {
		r.Response.Write(`test2`)
	})
	s.SetHandler(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set(ghttp.HeaderXUrlPath, "/test2")
		s.ServeHTTP(w, r)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(c.GetContent(ctx, "/"), "test2")
		t.Assert(c.GetContent(ctx, "/test/test"), "test2")
		t.Assert(c.GetContent(ctx, "/test1"), "test2")
	})
}
