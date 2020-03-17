// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func Test_Middleware_CORS(t *testing.T) {
	p := ports.PopRand()
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
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		// Common Checks.
		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/api.v2"), "Not Found")

		// GET request does not any route.
		resp, err := client.Get("/api.v2/user/list")
		gtest.Assert(err, nil)
		gtest.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		gtest.Assert(resp.StatusCode, 404)
		resp.Close()

		// POST request matches the route and CORS middleware.
		resp, err = client.Post("/api.v2/user/list")
		gtest.Assert(err, nil)
		gtest.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		gtest.Assert(resp.Header["Access-Control-Allow-Headers"][0], "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With")
		gtest.Assert(resp.Header["Access-Control-Allow-Methods"][0], "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE")
		gtest.Assert(resp.Header["Access-Control-Allow-Origin"][0], "*")
		gtest.Assert(resp.Header["Access-Control-Max-Age"][0], "3628800")
		resp.Close()
	})
	// OPTIONS GET
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "GET")
		resp, err := client.Options("/api.v2/user/list")
		gtest.Assert(err, nil)
		gtest.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 0)
		gtest.Assert(resp.ReadAllString(), "Not Found")
		gtest.Assert(resp.StatusCode, 404)
		resp.Close()
	})
	// OPTIONS POST
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		client.SetHeader("Access-Control-Request-Method", "POST")
		resp, err := client.Options("/api.v2/user/list")
		gtest.Assert(err, nil)
		gtest.Assert(len(resp.Header["Access-Control-Allow-Headers"]), 1)
		gtest.Assert(resp.StatusCode, 200)
		resp.Close()
	})
}
