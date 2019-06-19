// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 基本路由功能以及优先级测试
package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

// 基本路由功能测试
func Test_Router_DomainBasic(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/:name", func(r *ghttp.Request) {
		r.Response.Write("/:name")
	})
	d.BindHandler("/:name/update", func(r *ghttp.Request) {
		r.Response.Write(r.Get("name"))
	})
	d.BindHandler("/:name/:action", func(r *ghttp.Request) {
		r.Response.Write(r.Get("action"))
	})
	d.BindHandler("/:name/*any", func(r *ghttp.Request) {
		r.Response.Write(r.Get("any"))
	})
	d.BindHandler("/user/list/{field}.html", func(r *ghttp.Request) {
		r.Response.Write(r.Get("field"))
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))
		gtest.Assert(client.GetContent("/john"), "Not Found")
		gtest.Assert(client.GetContent("/john/update"), "Not Found")
		gtest.Assert(client.GetContent("/john/edit"), "Not Found")
		gtest.Assert(client.GetContent("/user/list/100.html"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))
		gtest.Assert(client.GetContent("/john"), "")
		gtest.Assert(client.GetContent("/john/update"), "john")
		gtest.Assert(client.GetContent("/john/edit"), "edit")
		gtest.Assert(client.GetContent("/user/list/100.html"), "100")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))
		gtest.Assert(client.GetContent("/john"), "")
		gtest.Assert(client.GetContent("/john/update"), "john")
		gtest.Assert(client.GetContent("/john/edit"), "edit")
		gtest.Assert(client.GetContent("/user/list/100.html"), "100")
	})
}

// 测试HTTP Method注册.
func Test_Router_DomainMethod(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("GET:/get", func(r *ghttp.Request) {

	})
	d.BindHandler("POST:/post", func(r *ghttp.Request) {

	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp1, err := client.Get("/get")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 404)

		resp2, err := client.Post("/get")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/post")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 404)

		resp4, err := client.Post("/post")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 404)
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		resp1, err := client.Get("/get")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 200)

		resp2, err := client.Post("/get")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/post")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 404)

		resp4, err := client.Post("/post")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 200)
	})

	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		resp1, err := client.Get("/get")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 200)

		resp2, err := client.Post("/get")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/post")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 404)

		resp4, err := client.Post("/post")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 200)
	})
}

// 测试状态返回.
func Test_Router_DomainStatus(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/200", func(r *ghttp.Request) {
		r.Response.WriteStatus(200)
	})
	d.BindHandler("/300", func(r *ghttp.Request) {
		r.Response.WriteStatus(300)
	})
	d.BindHandler("/400", func(r *ghttp.Request) {
		r.Response.WriteStatus(400)
	})
	d.BindHandler("/500", func(r *ghttp.Request) {
		r.Response.WriteStatus(500)
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		resp1, err := client.Get("/200")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 404)

		resp2, err := client.Get("/300")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 404)

		resp3, err := client.Get("/400")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 404)

		resp4, err := client.Get("/500")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 404)
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		resp1, err := client.Get("/200")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 200)

		resp2, err := client.Get("/300")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 300)

		resp3, err := client.Get("/400")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 400)

		resp4, err := client.Get("/500")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 500)
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		resp1, err := client.Get("/200")
		defer resp1.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp1.StatusCode, 200)

		resp2, err := client.Get("/300")
		defer resp2.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp2.StatusCode, 300)

		resp3, err := client.Get("/400")
		defer resp3.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp3.StatusCode, 400)

		resp4, err := client.Get("/500")
		defer resp4.Close()
		gtest.Assert(err, nil)
		gtest.Assert(resp4.StatusCode, 500)
	})
}

// 自定义状态码处理.
func Test_Router_DomainCustomStatusHandler(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	d.BindStatusHandler(404, func(r *ghttp.Request) {
		r.Response.Write("404 page")
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
		gtest.Assert(client.GetContent("/ThisDoesNotExist"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "hello")
		gtest.Assert(client.GetContent("/ThisDoesNotExist"), "404 page")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "hello")
		gtest.Assert(client.GetContent("/ThisDoesNotExist"), "404 page")
	})
}

// 测试不存在的路由.
func Test_Router_Domain404(t *testing.T) {
	p := ports.PopRand()
	s := g.Server(p)
	d := s.Domain("localhost, local")
	d.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello")
	})
	s.SetPort(p)
	s.SetDumpRouteMap(false)
	s.Start()
	defer s.Shutdown()

	// 等待启动完成
	time.Sleep(time.Second)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Not Found")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://localhost:%d", p))

		gtest.Assert(client.GetContent("/"), "hello")
	})
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://local:%d", p))

		gtest.Assert(client.GetContent("/"), "hello")
	})
}
