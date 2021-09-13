// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"testing"
	"time"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"github.com/gogf/gf/test/gtest"
)

func Test_ConfigFromMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.Map{
			"address":         ":8199",
			"readTimeout":     "60s",
			"indexFiles":      g.Slice{"index.php", "main.php"},
			"errorLogEnabled": true,
			"cookieMaxAge":    "1y",
		}
		config, err := ghttp.ConfigFromMap(m)
		t.Assert(err, nil)
		d1, _ := time.ParseDuration(gconv.String(m["readTimeout"]))
		d2, _ := time.ParseDuration(gconv.String(m["cookieMaxAge"]))
		t.Assert(config.Address, m["address"])
		t.Assert(config.ReadTimeout, d1)
		t.Assert(config.CookieMaxAge, d2)
		t.Assert(config.IndexFiles, m["indexFiles"])
		t.Assert(config.ErrorLogEnabled, m["errorLogEnabled"])
	})
}

func Test_SetConfigWithMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := g.Map{
			"Address": ":8199",
			//"ServerRoot":       "/var/www/MyServerRoot",
			"IndexFiles":       g.Slice{"index.php", "main.php"},
			"AccessLogEnabled": true,
			"ErrorLogEnabled":  true,
			"PProfEnabled":     true,
			"LogPath":          "/tmp/log/MyServerLog",
			"SessionIdName":    "MySessionId",
			"SessionPath":      "/tmp/MySessionStoragePath",
			"SessionMaxAge":    24 * time.Hour,
		}
		s := g.Server()
		err := s.SetConfigWithMap(m)
		t.Assert(err, nil)
	})
}

func Test_ClientMaxBodySize(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/", func(r *ghttp.Request) {
			r.Response.Write(r.GetBodyString())
		})
	})
	m := g.Map{
		"Address":           p,
		"ClientMaxBodySize": "1k",
	}
	gtest.Assert(s.SetConfigWithMap(m), nil)
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		data := make([]byte, 1056)
		for i := 0; i < 1056; i++ {
			data[i] = 'a'
		}
		t.Assert(
			gstr.Trim(c.PostContent("/", data)),
			data[:1024],
		)
	})
}

func Test_ClientMaxBodySize_File(t *testing.T) {
	p, _ := ports.PopRand()
	s := g.Server(p)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.POST("/", func(r *ghttp.Request) {
			r.GetUploadFile("file")
			r.Response.Write("ok")
		})
	})
	m := g.Map{
		"Address":           p,
		"ErrorLogEnabled":   false,
		"ClientMaxBodySize": "1k",
	}
	gtest.Assert(s.SetConfigWithMap(m), nil)
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	// ok
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		path := gfile.TempDir(gtime.TimestampNanoStr())
		data := make([]byte, 512)
		for i := 0; i < 512; i++ {
			data[i] = 'a'
		}
		t.Assert(gfile.PutBytes(path, data), nil)
		defer gfile.Remove(path)
		t.Assert(
			gstr.Trim(c.PostContent("/", "name=john&file=@file:"+path)),
			"ok",
		)
	})

	// too large
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		path := gfile.TempDir(gtime.TimestampNanoStr())
		data := make([]byte, 1056)
		for i := 0; i < 1056; i++ {
			data[i] = 'a'
		}
		t.Assert(gfile.PutBytes(path, data), nil)
		defer gfile.Remove(path)
		t.Assert(
			gstr.Trim(c.PostContent("/", "name=john&file=@file:"+path)),
			"Invalid Request: http: request body too large",
		)
	})
}
