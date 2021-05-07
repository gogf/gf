// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Log(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		logDir := gfile.TempDir(gtime.TimestampNanoStr())
		p, _ := ports.PopRand()
		s := g.Server(p)
		s.BindHandler("/hello", func(r *ghttp.Request) {
			r.Response.Write("hello")
		})
		s.BindHandler("/error", func(r *ghttp.Request) {
			panic("custom error")
		})
		s.SetLogPath(logDir)
		s.SetAccessLogEnabled(true)
		s.SetErrorLogEnabled(true)
		s.SetLogStdout(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		defer gfile.Remove(logDir)
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/hello"), "hello")
		t.Assert(client.GetContent("/error"), "custom error")

		logPath1 := gfile.Join(logDir, gtime.Now().Format("Y-m-d")+".log")
		t.Assert(gstr.Contains(gfile.GetContents(logPath1), "http server started listening on"), true)
		t.Assert(gstr.Contains(gfile.GetContents(logPath1), "HANDLER"), true)

		logPath2 := gfile.Join(logDir, "access-"+gtime.Now().Format("Ymd")+".log")
		//fmt.Println(gfile.GetContents(logPath2))
		t.Assert(gstr.Contains(gfile.GetContents(logPath2), " /hello "), true)

		logPath3 := gfile.Join(logDir, "error-"+gtime.Now().Format("Ymd")+".log")
		//fmt.Println(gfile.GetContents(logPath3))
		t.Assert(gstr.Contains(gfile.GetContents(logPath3), "custom error"), true)
	})
}
