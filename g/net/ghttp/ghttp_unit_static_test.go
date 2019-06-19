// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// 静态文件服务测试
package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
	"time"
)

func Test_Static_ServerRoot(t *testing.T) {
	// SetServerRoot with absolute path
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/index.htm", "index")
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "index")
		gtest.Assert(client.GetContent("/index.htm"), "index")
	})

	// SetServerRoot with relative path
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`static/test/%d`, p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/index.htm", "index")
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "index")
		gtest.Assert(client.GetContent("/index.htm"), "index")
	})
}

func Test_Static_Folder_Forbidden(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/test.html", "test")
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/index.html"), "Not Found")
		gtest.Assert(client.GetContent("/test.html"), "test")
	})
}

func Test_Static_IndexFolder(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/test.html", "test")
		s.SetIndexFolder(true)
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.AssertNE(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/index.html"), "Not Found")
		gtest.Assert(client.GetContent("/test.html"), "test")
	})
}

func Test_Static_IndexFiles1(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/index.html", "index")
		gfile.PutContents(path+"/test.html", "test")
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "index")
		gtest.Assert(client.GetContent("/index.html"), "index")
		gtest.Assert(client.GetContent("/test.html"), "test")
	})
}

func Test_Static_IndexFiles2(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/test.html", "test")
		s.SetIndexFiles([]string{"index.html", "test.html"})
		s.SetServerRoot(path)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "test")
		gtest.Assert(client.GetContent("/index.html"), "Not Found")
		gtest.Assert(client.GetContent("/test.html"), "test")
	})
}

func Test_Static_AddSearchPath1(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path1 := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		path2 := fmt.Sprintf(`%s/ghttp/static/test/%d/%d`, gfile.TempDir(), p, p)
		defer gfile.Remove(path1)
		defer gfile.Remove(path2)
		gfile.PutContents(path2+"/test.html", "test")
		s.SetServerRoot(path1)
		s.AddSearchPath(path2)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/test.html"), "test")
	})
}

func Test_Static_AddSearchPath2(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path1 := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		path2 := fmt.Sprintf(`%s/ghttp/static/test/%d/%d`, gfile.TempDir(), p, p)
		defer gfile.Remove(path1)
		defer gfile.Remove(path2)
		gfile.PutContents(path1+"/test.html", "test1")
		gfile.PutContents(path2+"/test.html", "test2")
		s.SetServerRoot(path1)
		s.AddSearchPath(path2)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/test.html"), "test1")
	})
}

func Test_Static_AddStaticPath(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path1 := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		path2 := fmt.Sprintf(`%s/ghttp/static/test/%d/%d`, gfile.TempDir(), p, p)
		defer gfile.Remove(path1)
		defer gfile.Remove(path2)
		gfile.PutContents(path1+"/test.html", "test1")
		gfile.PutContents(path2+"/test.html", "test2")
		s.SetServerRoot(path1)
		s.AddStaticPath("/my-test", path2)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/test.html"), "test1")
		gtest.Assert(client.GetContent("/my-test/test.html"), "test2")
	})
}

func Test_Static_AddStaticPath_Priority(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path1 := fmt.Sprintf(`%s/ghttp/static/test/%d/test`, gfile.TempDir(), p)
		path2 := fmt.Sprintf(`%s/ghttp/static/test/%d/%d/test`, gfile.TempDir(), p, p)
		defer gfile.Remove(path1)
		defer gfile.Remove(path2)
		gfile.PutContents(path1+"/test.html", "test1")
		gfile.PutContents(path2+"/test.html", "test2")
		s.SetServerRoot(path1)
		s.AddStaticPath("/test", path2)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/test.html"), "test1")
		gtest.Assert(client.GetContent("/test/test.html"), "test2")
	})
}

func Test_Static_Rewrite(t *testing.T) {
	gtest.Case(t, func() {
		p := ports.PopRand()
		s := g.Server(p)
		path := fmt.Sprintf(`%s/ghttp/static/test/%d`, gfile.TempDir(), p)
		defer gfile.Remove(path)
		gfile.PutContents(path+"/test1.html", "test1")
		gfile.PutContents(path+"/test2.html", "test2")
		s.SetServerRoot(path)
		s.SetRewrite("/test.html", "/test1.html")
		s.SetRewriteMap(g.MapStrStr{
			"/my-test1": "/test1.html",
			"/my-test2": "/test2.html",
		})
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(time.Second)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		gtest.Assert(client.GetContent("/"), "Forbidden")
		gtest.Assert(client.GetContent("/test.html"), "test1")
		gtest.Assert(client.GetContent("/test1.html"), "test1")
		gtest.Assert(client.GetContent("/test2.html"), "test2")
		gtest.Assert(client.GetContent("/my-test1"), "test1")
		gtest.Assert(client.GetContent("/my-test2"), "test2")
	})
}
