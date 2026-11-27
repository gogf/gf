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
	"github.com/gogf/gf/v2/net/gsvc"
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

		var logFilePath = fmt.Sprintf("/tmp/log/%s.log", gtime.Now().Format("Y-m-d"))
		defer gfile.RemoveFile(logFilePath)
		t.Assert(gfile.Exists(logFilePath), true)
		t.Assert(s.Logger().GetConfig().Prefix, "TestLogger")
	})

}

func Test_Grpcx_Grpc_Server_Config_Timeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := Server.NewConfig()
		t.AssertNil(c.SetWithMap(map[string]any{
			"timeout": "30s",
			"methodTimeouts": map[string]any{
				"/protobuf.Greeter/SayHello": "5s",
			},
		}))
		t.Assert(c.Timeout, 30*time.Second)
		t.Assert(c.MethodTimeouts["/protobuf.Greeter/SayHello"], 5*time.Second)
	})
}

func Test_Grpcx_Grpc_Server_InvalidServiceName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cfg := Server.NewConfig()
		cfg.Address = "127.0.0.1:0"
		cfg.Endpoints = []string{"127.0.0.1:8080"}
		s := Server.New(cfg)
		s.Service(&gsvc.LocalService{Name: "gogf-grpcx/demo2"})
		var recovered any
		func() {
			defer func() {
				recovered = recover()
			}()
			s.doServiceRegister()
		}()
		t.AssertNE(recovered, nil)
	})
}
