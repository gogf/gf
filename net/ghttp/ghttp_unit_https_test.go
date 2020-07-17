// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp_test

import (
	"fmt"
	"github.com/jin502437344/gf/debug/gdebug"
	"github.com/jin502437344/gf/os/gtime"
	"github.com/jin502437344/gf/text/gstr"
	"testing"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_HTTPS_Basic(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	s.EnableHTTPS(
		gdebug.TestDataPath("https", "server.crt"),
		gdebug.TestDataPath("https", "server.key"),
	)
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// HTTP
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		t.AssertIN(gstr.Trim(c.GetContent("/")), g.Slice{"", "Client sent an HTTP request to an HTTPS server."})
		t.AssertIN(gstr.Trim(c.GetContent("/test")), g.Slice{"", "Client sent an HTTP request to an HTTPS server."})
	})
	// HTTPS
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("https://127.0.0.1:%d", p))
		t.Assert(c.GetContent("/"), "Not Found")
		t.Assert(c.GetContent("/test"), "test")
	})
}

func Test_HTTPS_HTTP_Basic(t *testing.T) {
	var (
		portHttp, _  = ports.PopRand()
		portHttps, _ = ports.PopRand()
	)
	s := g.Server(gtime.TimestampNanoStr())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	s.EnableHTTPS(
		gdebug.TestDataPath("https", "server.crt"),
		gdebug.TestDataPath("https", "server.key"),
	)
	s.SetPort(portHttp)
	s.SetHTTPSPort(portHttps)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// HTTP
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", portHttp))
		t.Assert(c.GetContent("/"), "Not Found")
		t.Assert(c.GetContent("/test"), "test")
	})
	// HTTPS
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("https://127.0.0.1:%d", portHttps))
		t.Assert(c.GetContent("/"), "Not Found")
		t.Assert(c.GetContent("/test"), "test")
	})
}
