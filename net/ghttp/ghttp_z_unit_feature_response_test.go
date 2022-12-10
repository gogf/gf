// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
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
