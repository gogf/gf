// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/v2/encoding/gxml"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gview"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Response_ServeFile(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/ServeFile", func(r *ghttp.Request) {
		filePath := r.GetQuery("filePath")
		r.Response.ServeFile(filePath.String())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		srcPath := gtest.DataPath("upload", "file1.txt")
		t.Assert(client.GetContent(ctx, "/ServeFile", "filePath=file1.txt"), "Not Found")

		t.Assert(
			client.GetContent(ctx, "/ServeFile", "filePath="+srcPath),
			"file1.txt: This file is for uploading unit test case.")

		t.Assert(
			strings.Contains(
				client.GetContent(ctx, "/ServeFile", "filePath=files/server.key"),
				"BEGIN RSA PRIVATE KEY"),
			true)
	})
}

func Test_Response_ServeFileDownload(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/ServeFileDownload", func(r *ghttp.Request) {
		filePath := r.GetQuery("filePath")
		r.Response.ServeFileDownload(filePath.String())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		srcPath := gtest.DataPath("upload", "file1.txt")
		t.Assert(client.GetContent(ctx, "/ServeFileDownload", "filePath=file1.txt"), "Not Found")

		t.Assert(
			client.GetContent(ctx, "/ServeFileDownload", "filePath="+srcPath),
			"file1.txt: This file is for uploading unit test case.")

		t.Assert(
			strings.Contains(
				client.GetContent(ctx, "/ServeFileDownload", "filePath=files/server.key"),
				"BEGIN RSA PRIVATE KEY"),
			true)
	})
}

func Test_Response_Redirect(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("RedirectResult")
	})
	s.BindHandler("/RedirectTo", func(r *ghttp.Request) {
		r.Response.RedirectTo("/")
	})
	s.BindHandler("/RedirectTo301", func(r *ghttp.Request) {
		r.Response.RedirectTo("/", http.StatusMovedPermanently)
	})
	s.BindHandler("/RedirectBack", func(r *ghttp.Request) {
		r.Response.RedirectBack()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/RedirectTo"), "RedirectResult")
		t.Assert(client.GetContent(ctx, "/RedirectTo301"), "RedirectResult")
		t.Assert(client.SetHeader("Referer", "/").GetContent(ctx, "/RedirectBack"), "RedirectResult")
	})
}

func Test_Response_Buffer(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/Buffer", func(r *ghttp.Request) {
		name := r.GetQuery("name").Bytes()
		r.Response.SetBuffer(name)
		buffer := r.Response.Buffer()
		r.Response.ClearBuffer()
		r.Response.Write(buffer)
	})
	s.BindHandler("/BufferString", func(r *ghttp.Request) {
		name := r.GetQuery("name").Bytes()
		r.Response.SetBuffer(name)
		bufferString := r.Response.BufferString()
		r.Response.ClearBuffer()
		r.Response.Write(bufferString)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(ctx, "/Buffer", "name=john"), []byte("john"))
		t.Assert(client.GetContent(ctx, "/BufferString", "name=john"), "john")
	})
}

func Test_Response_WriteTpl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New(gtest.DataPath("template", "basic"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTpl("noexist.html", g.Map{
				"name": "john",
			})
			t.AssertNE(err, nil)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.AssertNE(client.GetContent(ctx, "/"), "Name:john")
	})
}

func Test_Response_WriteTplDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetDefaultFile(gtest.DataPath("template", "basic", "index.html"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTplDefault(g.Map{"name": "john"})
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Name:john")
	})
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetDefaultFile(gtest.DataPath("template", "basic", "noexit.html"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			err := r.Response.WriteTplDefault(g.Map{"name": "john"})
			t.AssertNil(err)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.AssertNE(client.GetContent(ctx, "/"), "Name:john")
	})
}

func Test_Response_ParseTplDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetDefaultFile(gtest.DataPath("template", "basic", "index.html"))
		s := g.Server(guid.S())
		s.SetView(v)
		s.BindHandler("/", func(r *ghttp.Request) {
			res, err := r.Response.ParseTplDefault(g.Map{"name": "john"})
			t.AssertNil(err)
			r.Response.Write(res)
		})
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "Name:john")
	})
}

func Test_Response_Write(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}
	s := g.Server(guid.S())
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write()
	})
	s.BindHandler("/WriteOverExit", func(r *ghttp.Request) {
		r.Response.Write("WriteOverExit")
		r.Response.WriteOverExit("")
	})
	s.BindHandler("/WritefExit", func(r *ghttp.Request) {
		r.Response.WritefExit("%s", "WritefExit")
	})
	s.BindHandler("/Writeln", func(r *ghttp.Request) {
		name := r.GetQuery("name")
		r.Response.Writeln(name)
	})
	s.BindHandler("/WritelnNil", func(r *ghttp.Request) {
		r.Response.Writeln()
	})
	s.BindHandler("/Writefln", func(r *ghttp.Request) {
		name := r.GetQuery("name")
		r.Response.Writefln("%s", name)
	})
	s.BindHandler("/WriteJson", func(r *ghttp.Request) {
		m := map[string]string{"name": "john"}
		if bytes, err := json.Marshal(m); err == nil {
			r.Response.WriteJson(bytes)
		}
	})
	s.BindHandler("/WriteJsonP", func(r *ghttp.Request) {
		m := map[string]string{"name": "john"}
		if bytes, err := json.Marshal(m); err == nil {
			r.Response.WriteJsonP(bytes)
		}
	})
	s.BindHandler("/WriteJsonPWithStruct", func(r *ghttp.Request) {
		user := User{"john"}
		r.Response.WriteJsonP(user)
	})
	s.BindHandler("/WriteXml", func(r *ghttp.Request) {
		m := map[string]interface{}{"name": "john"}
		if bytes, err := gxml.Encode(m); err == nil {
			r.Response.WriteXml(bytes)
		}
	})
	s.BindHandler("/WriteXmlWithStruct", func(r *ghttp.Request) {
		user := User{"john"}
		r.Response.WriteXml(user)
	})

	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(client.GetContent(ctx, "/"), "")
		t.Assert(client.GetContent(ctx, "/WriteOverExit"), "")
		t.Assert(client.GetContent(ctx, "/WritefExit"), "WritefExit")
		t.Assert(client.GetContent(ctx, "/Writeln"), "\n")
		t.Assert(client.GetContent(ctx, "/WritelnNil"), "\n")
		t.Assert(client.GetContent(ctx, "/Writeln", "name=john"), "john\n")
		t.Assert(client.GetContent(ctx, "/Writefln", "name=john"), "john\n")
		t.Assert(client.GetContent(ctx, "/WriteJson"), "{\"name\":\"john\"}")
		t.Assert(client.GetContent(ctx, "/WriteJsonP"), "{\"name\":\"john\"}")
		t.Assert(client.GetContent(ctx, "/WriteJsonPWithStruct"), "{\"name\":\"john\"}")
		t.Assert(client.GetContent(ctx, "/WriteJsonPWithStruct", "callback=callback"),
			"callback({\"name\":\"john\"})")
		t.Assert(client.GetContent(ctx, "/WriteXml"), "<name>john</name>")
		t.Assert(client.GetContent(ctx, "/WriteXmlWithStruct"), "<name>john</name>")
	})
}
