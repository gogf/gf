// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

// static service testing.

package ghttp_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	. "github.com/jin502437344/gf/test/gtest"
)

func TestServer_EnablePProf(t *testing.T) {
	C(t, func(t *T) {
		p, _ := ports.PopRand()
		s := g.Server(p)
		s.EnablePProf("/pprof")
		s.SetDumpRouterMap(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		r, err := client.Get("/pprof/index")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		r, err = client.Get("/pprof/cmdline")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		//r, err = client.Get("/pprof/profile")
		//Assert(err, nil)
		//Assert(r.StatusCode, 200)
		//r.Close()

		r, err = client.Get("/pprof/symbol")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()

		r, err = client.Get("/pprof/trace")
		Assert(err, nil)
		Assert(r.StatusCode, 200)
		r.Close()
	})

}
