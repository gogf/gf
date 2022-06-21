// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Client_Server_Tracing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := 8888
		s := g.Server(p)
		s.BindHandler("/", func(r *ghttp.Request) {
			ctx := r.Context()
			g.Log().Print(ctx, "GetTraceID:", gtrace.GetTraceID(ctx))
			r.Response.Write(gtrace.GetTraceID(ctx))
		})
		s.SetPort(p)
		s.SetDumpRouterMap(false)
		t.AssertNil(s.Start())
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		ctx := gctx.New()

		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)
		t.Assert(gtrace.IsUsingDefaultProvider(), true)
		t.Assert(client.GetContent(ctx, "/"), gtrace.GetTraceID(ctx))
		t.Assert(client.GetContent(ctx, "/"), gctx.CtxId(ctx))
	})
}

func Test_WithTraceID(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p := 8889
		s := g.Server(p)
		s.BindHandler("/", func(r *ghttp.Request) {
			ctx := r.Context()
			r.Response.Write(gtrace.GetTraceID(ctx))
		})
		s.SetPort(p)
		s.SetDumpRouterMap(false)
		t.AssertNil(s.Start())
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		ctx, err := gtrace.WithTraceID(context.TODO(), traceID.String())
		t.AssertNil(err)

		prefix := fmt.Sprintf("http://127.0.0.1:%d", p)
		client := g.Client()
		client.SetPrefix(prefix)
		t.Assert(gtrace.IsUsingDefaultProvider(), true)
		t.Assert(client.GetContent(ctx, "/"), gtrace.GetTraceID(ctx))
		t.Assert(client.GetContent(ctx, "/"), traceIDStr)
	})
}
