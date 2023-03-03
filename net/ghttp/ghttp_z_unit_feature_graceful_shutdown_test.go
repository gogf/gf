// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_GracefulShutdownWithSIGTERM(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/url", func(r *ghttp.Request) {
		time.Sleep(time.Second * 3)
		r.Response.Write(r.GetUrl())
	})
	s.SetDumpRouterMap(false)

	done := make(chan struct{})

	go func() {
		time.Sleep(100 * time.Millisecond)

		gtest.C(t, func(t *gtest.T) {
			prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

			client := g.Client()
			client.SetPrefix(prefix)

			// send signal SIGTERM after 200ms
			go func() {
				time.Sleep(200 * time.Millisecond)
				pid := os.Getpid()
				syscall.Kill(pid, syscall.SIGTERM)
			}()

			result := client.GetContent(ctx, "/url")
			expected := prefix + "/url"

			// execute `done <- struct{}{}` even if `t.AssertEQ(result, expected)` failed
			defer func() {
				done <- struct{}{}
			}()

			t.AssertEQ(result, expected)
		})
	}()

	s.Run()

	<-done
}

func Test_GracefulShutdownWithSIGINT(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/url", func(r *ghttp.Request) {
		time.Sleep(time.Second * 3)
		r.Response.Write(r.GetUrl())
	})
	s.SetDumpRouterMap(false)

	done := make(chan struct{})

	go func() {
		time.Sleep(100 * time.Millisecond)

		gtest.C(t, func(t *gtest.T) {
			prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

			client := g.Client()
			client.SetPrefix(prefix)

			// send signal SIGINT after 200ms
			go func() {
				time.Sleep(200 * time.Millisecond)
				pid := os.Getpid()
				syscall.Kill(pid, syscall.SIGINT)
			}()

			result := client.GetContent(ctx, "/url")
			expected := ""

			// execute `done <- struct{}{}` even if `t.AssertEQ(result, expected)` failed
			defer func() {
				done <- struct{}{}
			}()

			t.AssertEQ(result, expected)
		})
	}()

	s.Run()

	<-done
}
