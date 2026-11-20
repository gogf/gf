// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type closingReader struct {
	io.Reader
	closed chan struct{}
	once   sync.Once
}

func (c *closingReader) Close() error {
	c.once.Do(func() {
		close(c.closed)
	})
	return nil
}

func Test_Response_ServeReader(t *testing.T) {
	var (
		content = strings.Repeat("ServeReader content. ", 4096)
		closed  = make(chan struct{})
	)
	s := g.Server(guid.S())
	s.BindHandler("/known", func(r *ghttp.Request) {
		r.Response.ServeReader(strings.NewReader(content), int64(len(content)))
	})
	s.BindHandler("/unknown", func(r *ghttp.Request) {
		r.Response.ServeReader(strings.NewReader(content), -1)
	})
	s.BindHandler("/closer", func(r *ghttp.Request) {
		r.Response.ServeReader(&closingReader{Reader: strings.NewReader(content), closed: closed}, int64(len(content)))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		resp, err := client.Get(ctx, "/known")
		t.AssertNil(err)
		t.Assert(resp.ReadAllString(), content)
		t.Assert(resp.ContentLength, int64(len(content)))
		t.Assert(resp.Header.Get("Accept-Ranges"), "")

		resp, err = client.Get(ctx, "/unknown")
		t.AssertNil(err)
		t.Assert(resp.ReadAllString(), content)
		t.Assert(resp.ContentLength, int64(-1))
		t.Assert(resp.Header.Get("Accept-Ranges"), "")

		resp, err = client.Head(ctx, "/known")
		t.AssertNil(err)
		t.Assert(resp.ContentLength, int64(len(content)))
		t.Assert(resp.ReadAllString(), "")

		resp, err = client.Head(ctx, "/unknown")
		t.AssertNil(err)
		t.Assert(resp.ContentLength, int64(-1))
		t.Assert(resp.ReadAllString(), "")

		resp, err = client.Get(ctx, "/closer")
		t.AssertNil(err)
		t.Assert(resp.ReadAllString(), content)
		select {
		case <-closed:
		case <-time.After(time.Second):
			t.Error("reader was not closed after the content was sent")
		}
	})
}

func Test_Response_ServeReaderDownload(t *testing.T) {
	var content = strings.Repeat("ServeReaderDownload content. ", 2048)
	s := g.Server(guid.S())
	s.BindHandler("/download", func(r *ghttp.Request) {
		if contentType := r.GetQuery("type").String(); contentType != "" {
			r.Response.ServeReaderDownload(strings.NewReader(content), int64(len(content)), r.GetQuery("name").String(), contentType)
			return
		}
		r.Response.ServeReaderDownload(strings.NewReader(content), int64(len(content)), r.GetQuery("name").String())
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		resp, err := client.Get(ctx, "/download", g.Map{"name": "report.pdf"})
		t.AssertNil(err)
		t.Assert(resp.ReadAllString(), content)
		t.Assert(resp.ContentLength, int64(len(content)))
		t.Assert(resp.Header.Get("Content-Disposition"), "attachment; filename=report.pdf")
		t.Assert(resp.Header.Get("Content-Type"), "application/octet-stream")
		t.Assert(resp.Header.Get("Access-Control-Expose-Headers"), "Content-Disposition")
		t.Assert(resp.Header.Get("Accept-Ranges"), "")

		resp, err = client.Get(ctx, "/download", g.Map{"name": "report.pdf", "type": "application/pdf"})
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Type"), "application/pdf")

		resp, err = client.Get(ctx, "/download", g.Map{"name": "报表 2026.pdf"})
		t.AssertNil(err)
		t.Assert(resp.Header.Get("Content-Disposition"), `attachment; filename*=utf-8''%E6%8A%A5%E8%A1%A8%202026.pdf`)

		resp, err = client.Head(ctx, "/download?name=report.pdf")
		t.AssertNil(err)
		t.Assert(resp.ContentLength, int64(len(content)))
		t.Assert(resp.Header.Get("Content-Disposition"), "attachment; filename=report.pdf")
		t.Assert(resp.ReadAllString(), "")
	})
}
