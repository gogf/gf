// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_MiddlewareGzip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server()
		s.Use(ghttp.MiddlewareGzip)

		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("Hello, World!")
		})

		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		resp, err := c.Get(context.TODO(), "/")
		t.AssertNil(err)
		defer resp.Close()

		t.Assert(resp.Header.Get("Content-Encoding"), "gzip")

		gzipReader, err := gzip.NewReader(resp.Body)
		t.AssertNil(err)
		defer gzipReader.Close()

		body, err := io.ReadAll(gzipReader)
		fmt.Printf("t: %v\n", body)
		t.AssertNil(err)
		t.Assert(string(body), "Hello, World!")
	})
}
