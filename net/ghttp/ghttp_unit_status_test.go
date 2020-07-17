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
	"github.com/jin502437344/gf/test/gtest"
)

func Test_StatusHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		p, _ := ports.PopRand()
		s := g.Server(p)
		s.BindStatusHandlerByMap(map[int]ghttp.HandlerFunc{
			404: func(r *ghttp.Request) { r.Response.WriteOver("404") },
			502: func(r *ghttp.Request) { r.Response.WriteOver("502") },
		})
		s.BindHandler("/502", func(r *ghttp.Request) {
			r.Response.WriteStatusExit(502)
		})
		s.SetDumpRouterMap(false)
		s.SetPort(p)
		s.Start()
		defer s.Shutdown()
		time.Sleep(100 * time.Millisecond)
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		t.Assert(client.GetContent("/404"), "404")
		t.Assert(client.GetContent("/502"), "502")
	})
}
