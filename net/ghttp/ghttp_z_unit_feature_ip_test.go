// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestRequest_GetRemoteIp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write(r.GetRemoteIp())
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		clientV4 := g.Client()
		clientV4.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		clientV6 := g.Client()
		clientV6.SetPrefix(fmt.Sprintf("http://[::1]:%d", s.GetListenedPort()))

		t.Assert(clientV4.GetContent(ctx, "/"), "127.0.0.1")
		t.Assert(clientV6.GetContent(ctx, "/"), "::1")
	})
}
