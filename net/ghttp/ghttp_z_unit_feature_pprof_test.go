// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

func TestServer_EnablePProf(t *testing.T) {
	C(t, func(t *T) {
		s := g.Server(guid.S())
		s.EnablePProf("/pprof")
		s.SetDumpRouterMap(false)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		urlPaths := []string{
			"/pprof/index", "/pprof/cmdline", "/pprof/symbol", "/pprof/trace",
		}
		for _, urlPath := range urlPaths {
			r, err := client.Get(ctx, urlPath)
			AssertNil(err)
			Assert(r.StatusCode, 200)
			AssertNil(r.Close())
		}
	})
}

func TestServer_StartPProfServer(t *testing.T) {
	C(t, func(t *T) {
		s, err := ghttp.StartPProfServer(":0")
		t.AssertNil(err)

		defer ghttp.ShutdownAllServer(ctx)

		time.Sleep(100 * time.Millisecond)
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d/debug", s.GetListenedPort()))

		urlPaths := []string{
			"/pprof/index", "/pprof/cmdline", "/pprof/symbol", "/pprof/trace",
		}
		for _, urlPath := range urlPaths {
			r, err := client.Get(ctx, urlPath)
			AssertNil(err)
			Assert(r.StatusCode, 200)
			AssertNil(r.Close())
		}
	})
}
