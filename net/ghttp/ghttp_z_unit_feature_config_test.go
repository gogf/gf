// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_ConfigFromMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.Map{
			"address":         ":12345",
			"listeners":       nil,
			"readTimeout":     "60s",
			"indexFiles":      g.Slice{"index.php", "main.php"},
			"errorLogEnabled": true,
			"cookieMaxAge":    "1y",
			"cookieSameSite":  "lax",
			"cookieSecure":    true,
			"cookieHttpOnly":  true,
		}
		config, err := ghttp.ConfigFromMap(m)
		t.AssertNil(err)
		d1, _ := time.ParseDuration(gconv.String(m["readTimeout"]))
		d2, _ := time.ParseDuration(gconv.String(m["cookieMaxAge"]))
		t.Assert(config.Address, m["address"])
		t.Assert(config.ReadTimeout, d1)
		t.Assert(config.CookieMaxAge, d2)
		t.Assert(config.IndexFiles, m["indexFiles"])
		t.Assert(config.ErrorLogEnabled, m["errorLogEnabled"])
		t.Assert(config.CookieSameSite, m["cookieSameSite"])
		t.Assert(config.CookieSecure, m["cookieSecure"])
		t.Assert(config.CookieHttpOnly, m["cookieHttpOnly"])
	})
}

func Test_SetConfigWithMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.Map{
			"Address": ":8199",
			// "ServerRoot":       "/var/www/MyServerRoot",
			"IndexFiles":       g.Slice{"index.php", "main.php"},
			"AccessLogEnabled": true,
			"ErrorLogEnabled":  true,
			"PProfEnabled":     true,
			"LogPath":          "/tmp/log/MyServerLog",
			"SessionIdName":    "MySessionId",
			"SessionPath":      "/tmp/MySessionStoragePath",
			"SessionMaxAge":    24 * time.Hour,
			"cookieSameSite":   "lax",
			"cookieSecure":     true,
			"cookieHttpOnly":   true,
		}
		s := g.Server()
		err := s.SetConfigWithMap(m)
		t.AssertNil(err)
	})
}

func Test_ClientMaxBodySize(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/", func(r *ghttp.Request) {
			r.Response.Write(r.GetBodyString())
		})
	})
	m := g.Map{
		"ClientMaxBodySize": "1k",
	}
	gtest.Assert(s.SetConfigWithMap(m), nil)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		data := make([]byte, 1056)
		for i := 0; i < 1056; i++ {
			data[i] = 'a'
		}
		t.Assert(
			gstr.Trim(c.PostContent(ctx, "/", data)),
			`Read from request Body failed: http: request body too large`,
		)
	})
}

func Test_ClientMaxBodySize_File(t *testing.T) {
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/", func(r *ghttp.Request) {
			r.GetUploadFile("file")
			r.Response.Write("ok")
		})
	})
	m := g.Map{
		"ErrorLogEnabled":   false,
		"ClientMaxBodySize": "1k",
	}
	gtest.Assert(s.SetConfigWithMap(m), nil)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// ok
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		path := gfile.Temp(gtime.TimestampNanoStr())
		data := make([]byte, 512)
		for i := 0; i < 512; i++ {
			data[i] = 'a'
		}
		t.Assert(gfile.PutBytes(path, data), nil)
		defer gfile.Remove(path)
		t.Assert(
			gstr.Trim(c.PostContent(ctx, "/", "name=john&file=@file:"+path)),
			"ok",
		)
	})

	// too large
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		path := gfile.Temp(gtime.TimestampNanoStr())
		data := make([]byte, 1056)
		for i := 0; i < 1056; i++ {
			data[i] = 'a'
		}
		t.Assert(gfile.PutBytes(path, data), nil)
		defer gfile.Remove(path)
		t.Assert(
			true,
			strings.Contains(
				gstr.Trim(c.PostContent(ctx, "/", "name=john&file=@file:"+path)),
				"http: request body too large",
			),
		)
	})
}
