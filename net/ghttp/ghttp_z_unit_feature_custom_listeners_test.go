// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_SetSingleCustomListener(t *testing.T) {
	ln1, _ := net.Listen("tcp", ":0")
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})
	err := s.SetListener(ln1)
	gtest.AssertNil(err)

	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		t.Assert(
			gstr.Trim(c.GetContent(ctx, "/test")),
			"test",
		)
	})
}

func Test_SetMultipleCustomListeners(t *testing.T) {
	ln1, _ := net.Listen("tcp", ":0")
	ln2, _ := net.Listen("tcp", ":0")
	s := g.Server(guid.S())
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})
	})

	err := s.SetListener(ln1, ln2)
	gtest.AssertNil(err)

	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		ports := s.GetListenedPorts()
		t.Assert(len(ports), 2)

		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", ports[0]))

		t.Assert(
			gstr.Trim(c.GetContent(ctx, "/test")),
			"test",
		)

		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", ports[1]))

		t.Assert(
			gstr.Trim(c.GetContent(ctx, "/test")),
			"test",
		)
	})
}

func Test_SetWrongCustomListeners(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server(guid.S())
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/test", func(r *ghttp.Request) {
				r.Response.Write("test")
			})
		})
		err := s.SetListener(nil)
		t.AssertNQ(err, nil)
	})
}
