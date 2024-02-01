// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package file_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/guid"
)

var ctx = gctx.GetInitCtx()

func Test_HTTP_Registry(t *testing.T) {
	var (
		svcName = guid.S()
		dirPath = gfile.Temp(svcName)
	)
	defer gfile.Remove(dirPath)
	gsvc.SetRegistry(file.New(dirPath))

	s := g.Server(svcName)
	s.BindHandler("/http-registry", func(r *ghttp.Request) {
		r.Response.Write(svcName)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://%s", svcName))
		// GET
		t.Assert(client.GetContent(ctx, "/http-registry"), svcName)
	})
}

func Test_HTTP_Discovery_Disable(t *testing.T) {
	var (
		svcName = guid.S()
		dirPath = gfile.Temp(svcName)
	)
	defer gfile.Remove(dirPath)
	gsvc.SetRegistry(file.New(dirPath))

	s := g.Server(svcName)
	s.BindHandler("/http-registry", func(r *ghttp.Request) {
		r.Response.Write(svcName)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://%s", svcName))
		result, err := client.Get(ctx, "/http-registry")
		defer result.Close()
		t.AssertNil(err)
		t.Assert(result.ReadAllString(), svcName)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://%s", svcName))
		result, err := client.Discovery(nil).Get(ctx, "/http-registry")
		defer result.Close()
		t.AssertNE(err, nil)
	})
}

func Test_HTTP_Server_Endpoints(t *testing.T) {
	var (
		svcName = guid.S()
		dirPath = gfile.Temp(svcName)
	)
	defer gfile.Remove(dirPath)
	gsvc.SetRegistry(file.New(dirPath))

	endpoints := []string{"10.0.0.1:8000", "10.0.0.2:8000"}
	s := g.Server(svcName)
	s.SetEndpoints(endpoints)
	s.BindHandler("/http-registry", func(r *ghttp.Request) {
		r.Response.Write(svcName)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		service, err := gsvc.Get(ctx, svcName)
		t.AssertNil(err)
		t.Assert(service.GetEndpoints(), gstr.Join(endpoints, ","))
	})
}
