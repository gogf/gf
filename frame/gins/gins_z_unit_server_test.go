// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Server(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path                = gcfg.DefaultConfigFileName
			serverConfigContent = gtest.DataContent("server", "config.yaml")
			err                 = gfile.PutContents(path, serverConfigContent)
		)
		t.AssertNil(err)
		defer gfile.Remove(path)

		time.Sleep(time.Second)

		localInstances.Clear()
		defer localInstances.Clear()

		s := Server("tempByInstanceName")
		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("hello")
		})
		s.SetDumpRouterMap(false)
		t.AssertNil(s.Start())
		defer t.AssertNil(s.Shutdown())

		content := HttpClient().GetContent(gctx.New(), `http://127.0.0.1:8003/`)
		t.Assert(content, `hello`)
	})
}
