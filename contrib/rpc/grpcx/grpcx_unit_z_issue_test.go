// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx_test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var ctx = context.Background()

// https://github.com/gogf/gf/issues/3292
func Test_Issue3292(t *testing.T) {
	var (
		_ = grpcx.Client.MustNewGrpcClientConn(
			"127.0.0.1:8888",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
	)

	s := g.Server(guid.S())
	s.BindHandler("/url", func(r *ghttp.Request) {
		r.Response.Write(1)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		res, err := client.Get(ctx, "/url")
		t.AssertNil(err)
		defer res.Close()

		t.Assert(res.ReadAllString(), "1")
	})
}
