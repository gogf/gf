// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"net"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Test_SetRightCustomListeners(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server()
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/test", func(r *ghttp.Request) {
				r.Response.Write("test")
			})
		})
		s.SetAddr(":8199")
		ln, err := net.Listen("tcp", ":8199")
		t.AssertNil(err)
		s.SetListeners(map[int]net.Listener{8199: ln})

		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		s.GetListenedPort()

		t.AssertEQ(s.GetListenedPort(), 8199)
	})
}

func Test_SetWrongCustomListeners(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server()
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/test", func(r *ghttp.Request) {
				r.Response.Write("test")
			})
		})
		s.SetAddr(":8199")
		ln, err := net.Listen("tcp", ":8299")
		t.AssertNil(err)
		s.SetListeners(map[int]net.Listener{8199: ln})

		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		s.GetListenedPort()

		t.AssertEQ(s.GetListenedPort(), 8199)
	})
}
