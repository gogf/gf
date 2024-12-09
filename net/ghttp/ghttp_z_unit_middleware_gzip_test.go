// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Middleware_Gzip(t *testing.T) {
	s := g.Server(guid.S())
	// Routes with GZIP enabled
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(ghttp.MiddlewareGzip)
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write(strings.Repeat("Hello World! ", 1000))
		})
		group.ALL("/small", func(r *ghttp.Request) {
			r.Response.Write("Small response")
		})
	})

	// Routes without GZIP
	s.Group("/no-gzip", func(group *ghttp.RouterGroup) {
		group.ALL("/", func(r *ghttp.Request) {
			r.Response.Write(strings.Repeat("Hello World! ", 1000))
		})
	})

	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		// Test 1: Route with GZIP, client supports GZIP
		resp, err := client.Header(map[string]string{
			"Accept-Encoding": "gzip",
		}).Get(ctx, "/")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Encoding"), "gzip")

		reader, err := gzip.NewReader(resp.Body)
		t.AssertNil(err)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		t.AssertNil(err)
		expected := strings.Repeat("Hello World! ", 1000)
		t.Assert(len(content), len(expected))
		t.Assert(string(content), expected)

		// Test 2: Route with GZIP, client doesn't support GZIP
		resp, err = client.Header(map[string]string{}).Get(ctx, "/")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Encoding"), "")
		content, err = io.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(len(content), len(expected))
		t.Assert(string(content), expected)

		// Test 3: Route with GZIP, response too small
		resp, err = client.Header(map[string]string{
			"Accept-Encoding": "gzip",
		}).Get(ctx, "/small")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Encoding"), "")
		content, err = io.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(content), "Small response")

		// Test 4: Route without GZIP
		resp, err = client.Header(map[string]string{
			"Accept-Encoding": "gzip",
		}).Get(ctx, "/no-gzip/")
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Encoding"), "")
		content, err = io.ReadAll(resp.Body)
		t.AssertNil(err)
		t.Assert(string(content), expected)
	})
}
