// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Grpcx_Grpc_Server(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := Server.New()
		s.Start()
		time.Sleep(time.Millisecond * 100)
		defer s.Stop()
		s.serviceMu.Lock()
		defer s.serviceMu.Unlock()
		t.Assert(len(s.services) != 0, true)
	})
}

func Test_Grpcx_Grpc_Server_Address(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := Server.NewConfig()
		c.Address = "127.0.0.1:0"
		s := Server.New(c)
		s.Start()
		time.Sleep(time.Millisecond * 100)
		defer s.Stop()

		s.serviceMu.Lock()
		defer s.serviceMu.Unlock()
		t.Assert(len(s.services) != 0, true)
		t.Assert(gstr.Contains(s.services[0].GetEndpoints().String(), "127.0.0.1:"), true)
	})
}

func Test_Grpcx_Grpc_Server_Config(t *testing.T) {
	cfg := Server.NewConfig()
	addr := "10.0.0.29:80"
	cfg.Endpoints = []string{
		addr,
	}
	// cfg set one endpoint
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
	// cfg set more endpoints
	addr = "10.0.0.29:80,10.0.0.29:81"
	cfg.Endpoints = []string{
		"10.0.0.29:80",
		"10.0.0.29:81",
	}
	gtest.C(t, func(t *gtest.T) {
		s := Server.New(cfg)
		s.doServiceRegister()
		for _, svc := range s.services {
			t.Assert(svc.GetEndpoints().String(), addr)
		}
	})
}

func Test_Grpcx_Grpc_Server_Config_Logger(t *testing.T) {
	var (
		pwd       = gfile.Pwd()
		configDir = gfile.Join(gdebug.CallerDirectory(), "testdata", "configuration")
	)
	gtest.C(t, func(t *gtest.T) {
		err := gfile.Chdir(configDir)
		t.AssertNil(err)
		defer gfile.Chdir(pwd)

		s := Server.New()
		s.Start()
		time.Sleep(time.Millisecond * 100)
		defer s.Stop()

		var (
			logFilePath    = fmt.Sprintf("/tmp/log/%s.log", gtime.Now().Format("Y-m-d"))
			logFileContent = gfile.GetContents(logFilePath)
		)
		defer gfile.Remove(logFilePath)
		t.Assert(gfile.Exists(logFilePath), true)
		t.Assert(gstr.Contains(logFileContent, "TestLogger "), true)
	})

}
