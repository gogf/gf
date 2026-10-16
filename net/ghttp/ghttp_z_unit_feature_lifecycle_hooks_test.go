// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Server_Lifecycle_BeforeStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			hookExecuted = false
		)
		s := g.Server(gtest.DataPath("lifecycle-before-start"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetBeforeStart(func(s *ghttp.Server) error {
			hookExecuted = true
			t.Assert(s.GetListenedPort(), -1) // Not started yet
			return nil
		})

		err := s.Start()
		t.AssertNil(err)
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		t.Assert(hookExecuted, true)
	})
}

func Test_Server_Lifecycle_BeforeStart_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			serverName = gtest.DataPath("lifecycle-before-start-error")
		)
		s := g.Server(serverName)
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetBeforeStart(func(s *ghttp.Server) error {
			return fmt.Errorf("before start error")
		})

		err := s.Start()
		t.AssertNE(err, nil)
		t.Assert(ghttp.GetServer(serverName).Status(), ghttp.ServerStatusStopped)
	})
}

func Test_Server_Lifecycle_AfterStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			hookExecuted = false
			listenedPort = 0
		)
		s := g.Server(gtest.DataPath("lifecycle-after-start"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetAfterStart(func(s *ghttp.Server) error {
			hookExecuted = true
			listenedPort = s.GetListenedPort()
			t.Assert(listenedPort > 0, true) // Already started, port assigned
			return nil
		})

		err := s.Start()
		t.AssertNil(err)
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		t.Assert(hookExecuted, true)
		t.Assert(listenedPort > 0, true)
	})
}

func Test_Server_Lifecycle_MultipleHooks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			beforeCount = 0
			afterCount  = 0
			mu          sync.Mutex
		)
		s := g.Server(gtest.DataPath("lifecycle-multiple-hooks"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		// Register multiple before-start hooks
		s.SetBeforeStart(func(s *ghttp.Server) error {
			mu.Lock()
			beforeCount++
			mu.Unlock()
			return nil
		})
		s.SetBeforeStart(func(s *ghttp.Server) error {
			mu.Lock()
			beforeCount++
			mu.Unlock()
			return nil
		})

		// Register multiple after-start hooks
		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			afterCount++
			mu.Unlock()
			return nil
		})
		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			afterCount++
			mu.Unlock()
			return nil
		})

		err := s.Start()
		t.AssertNil(err)
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		t.Assert(beforeCount, 2)
		t.Assert(afterCount, 2)
	})
}

func Test_Server_Lifecycle_HooksOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			executionOrder []string
			mu             sync.Mutex
		)
		s := g.Server(gtest.DataPath("lifecycle-hooks-order"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetBeforeStart(func(s *ghttp.Server) error {
			mu.Lock()
			executionOrder = append(executionOrder, "before-1")
			mu.Unlock()
			return nil
		})
		s.SetBeforeStart(func(s *ghttp.Server) error {
			mu.Lock()
			executionOrder = append(executionOrder, "before-2")
			mu.Unlock()
			return nil
		})
		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			executionOrder = append(executionOrder, "after-1")
			mu.Unlock()
			return nil
		})
		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			executionOrder = append(executionOrder, "after-2")
			mu.Unlock()
			return nil
		})

		err := s.Start()
		t.AssertNil(err)
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		t.Assert(len(executionOrder), 4)
		t.Assert(executionOrder[0], "before-1")
		t.Assert(executionOrder[1], "before-2")
		t.Assert(executionOrder[2], "after-1")
		t.Assert(executionOrder[3], "after-2")
	})
}

func Test_Server_Lifecycle_BeforeStart_StopOnError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			secondHookExecuted = false
		)
		s := g.Server(gtest.DataPath("lifecycle-before-start-stop"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetBeforeStart(func(s *ghttp.Server) error {
			return fmt.Errorf("first hook error")
		})
		s.SetBeforeStart(func(s *ghttp.Server) error {
			secondHookExecuted = true
			return nil
		})

		err := s.Start()
		t.AssertNE(err, nil)
		t.Assert(secondHookExecuted, false) // Should not execute after error
	})
}

func Test_Server_Lifecycle_AfterStart_WithError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			firstHookExecuted  = false
			secondHookExecuted = false
			mu                 sync.Mutex
		)
		s := g.Server(gtest.DataPath("lifecycle-after-start-error"))
		s.BindHandler("/test", func(r *ghttp.Request) {
			r.Response.Write("test")
		})

		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			firstHookExecuted = true
			mu.Unlock()
			return fmt.Errorf("first hook error")
		})
		s.SetAfterStart(func(s *ghttp.Server) error {
			mu.Lock()
			secondHookExecuted = true
			mu.Unlock()
			return nil
		})

		err := s.Start()
		t.AssertNil(err) // Server should start successfully
		defer s.Shutdown()

		time.Sleep(100 * time.Millisecond)
		t.Assert(firstHookExecuted, true)
		t.Assert(secondHookExecuted, true) // Should continue even after error
	})
}
