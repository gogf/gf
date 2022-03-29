// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/test/gtest"
	"net"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func Test_SetSingleCustomListener(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p, _ := gtcp.GetFreePort()
		addr := fmt.Sprintf(":%d", p)
		s := g.Server(g.Map{
			"address": addr,
		})
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/test", func(r *ghttp.Request) {
				r.Response.Write("test")
			})
		})
		ln, err := net.Listen("tcp", addr)
		t.AssertNil(err)
		err = s.SetListener(ln)
		t.AssertNil(err)

		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
	})
}

func Test_SetMultipleCustomListeners(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Server()
		s.Group("/", func(group *ghttp.RouterGroup) {
			group.GET("/test", func(r *ghttp.Request) {
				r.Response.Write("test")
			})
		})
		p1, _ := gtcp.GetFreePort()
		p2, _ := gtcp.GetFreePort()

		ln1, err := net.Listen("tcp", fmt.Sprintf(":%d", p1))
		ln2, err := net.Listen("tcp", fmt.Sprintf(":%d", p2))
		err = s.SetListener(ln1, ln2)
		t.AssertEQ(err, nil)

		s.Start()
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		ports := []int{p1, p2}
		for _, p := range s.GetListenedPorts() {
			t.AssertIN(p, ports)
		}
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
		err := s.SetListener(nil)
		t.AssertNQ(err, nil)
	})
}
