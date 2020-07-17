// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package ghttp_test

import (
	"fmt"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"github.com/jin502437344/gf/test/gtest"
	"testing"
	"time"
)

func Test_Middleware_CORS1(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		group.POST("/user/list", func(r *ghttp.Request) {
			r.Response.Write("list")
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
		// Common Checks.
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/api.v2"), "Not Found")

		// GET request does not any route.
		resp, err := client.Get("/api.v2/user/list")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		t.Assert(resp.StatusCode, 404)
		resp.Close()

		// POST request matches the route and CORS middleware.
		resp, err = client.Post("/api.v2/user/list")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		t.Assert(resp.Header["Access-Control-Allow-Headers"][0], "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		t.Assert(resp.Header["Access-Control-Allow-Methods"][0], "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		t.Assert(resp.Header["Access-Control-Allow-Origin"][0], "*")
		t.Assert(resp.Header["Access-Control-Max-Age"][0], "3628800")
		resp.Close()
	})
	// OPTIONS GET
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "GET")
		resp, err := client.Options("/api.v2/user/list")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		t.Assert(resp.ReadAllString(), "Not Found")
		t.Assert(resp.StatusCode, 404)
		resp.Close()
	})
	// OPTIONS POST
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "POST")
		resp, err := client.Options("/api.v2/user/list")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		t.Assert(resp.StatusCode, 200)
		resp.Close()
	})
}

func Test_Middleware_CORS2(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/api.v2", func(group *ghttp.RouterGroup) {
		group.Middleware(MiddlewareCORS)
		group.GET("/user/list/{type}", func(r *ghttp.Request) {
			r.Response.Write(r.Get("type"))
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
		// Common Checks.
		t.Assert(client.GetContent("/"), "Not Found")
		t.Assert(client.GetContent("/api.v2"), "Not Found")
		// Get request.
		resp, err := client.Get("/api.v2/user/list/1")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		t.Assert(resp.Header["Access-Control-Allow-Headers"][0], "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		t.Assert(resp.Header["Access-Control-Allow-Methods"][0], "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		t.Assert(resp.Header["Access-Control-Allow-Origin"][0], "*")
		t.Assert(resp.Header["Access-Control-Max-Age"][0], "3628800")
		t.Assert(resp.ReadAllString(), "1")
		resp.Close()
	})
	// OPTIONS GET None.
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "GET")
		resp, err := client.Options("/api.v2/user")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		t.Assert(resp.StatusCode, 404)
		resp.Close()
	})
	// OPTIONS GET
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "GET")
		resp, err := client.Options("/api.v2/user/list/1")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		t.Assert(resp.StatusCode, 200)
		resp.Close()
	})
	// OPTIONS POST
	gtest.C(t, func(t *gtest.T) {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "POST")
		resp, err := client.Options("/api.v2/user/list/1")
		t.Assert(err, nil)
		t.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		t.Assert(resp.StatusCode, 404)
		resp.Close()
	})
}
