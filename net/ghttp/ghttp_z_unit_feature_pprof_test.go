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

	"github.com/gogf/gf/v2/frame/g"
	. "github.com/gogf/gf/v2/test/gtest"
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

		r, err := client.Get(ctx, "/pprof/index")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		r, err = client.Get(ctx, "/pprof/cmdline")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		//r, err = client.Get(ctx, "/pprof/profile")
		//Assert(err, nil)
		//Assert(r.StatusCode, 200)
		//r.Close()

		r, err = client.Get(ctx, "/pprof/symbol")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		r, err = client.Get(ctx, "/pprof/trace")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()
	})

}
