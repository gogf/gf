// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"net/url"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Static_IndexFolder_Escape(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx  = context.Background()
			root = t.TempDir()
		)
		t.AssertNil(gfile.PutContents(gfile.Join(root, "a&b #c.txt"), "test"))
		weirdDir := ""
		if runtime.GOOS != "windows" {
			weirdDir = `<img src=x onerror=alert(1)>`
			t.AssertNil(gfile.Mkdir(gfile.Join(root, weirdDir)))
			t.AssertNil(gfile.PutContents(gfile.Join(root, weirdDir, `x"y.txt`), "test"))
		}

		s := g.Server(guid.S())
		s.SetServerRoot(root)
		s.SetIndexFolder(true)
		s.SetDumpRouterMap(false)
		t.AssertNil(s.Start())
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		body := client.GetContent(ctx, "/")
		t.Assert(strings.Contains(body, `href="/a&amp;b%20%23c.txt"`), true)
		t.Assert(strings.Contains(body, `a&amp;b #c.txt`), true)

		if runtime.GOOS != "windows" {
			body = client.GetContent(ctx, "/"+url.PathEscape(weirdDir))
			t.Assert(strings.Contains(body, `Index of /&lt;img src=x onerror=alert(1)&gt;`), true)
			t.Assert(strings.Contains(body, `<img src=x onerror=alert(1)>`), false)
			t.Assert(strings.Contains(body, `href="/%3Cimg%20src=x%20onerror=alert%281%29%3E/x%22y.txt"`), true)
		}
	})
}
