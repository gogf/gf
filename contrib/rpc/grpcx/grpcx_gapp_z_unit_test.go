// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Unit tests for gapp.Server adapter wrapping GrpcServer.
package grpcx_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/testdata/controller"
)

func Test_Grpcx_GappServerAdapter_ImplementsGappServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := grpcx.Server.NewConfig()
		c.Name = guid.S()
		s := grpcx.Server.New(c)
		var _ gapp.Server = grpcx.NewGappServerAdapter(s)
	})
}

func Test_Grpcx_StartManagedDoesNotRegisterSignalHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := grpcx.Server.NewConfig()
		c.Name = guid.S()
		s := grpcx.Server.New(c)
		controller.Register(s)

		err := s.StartManaged()
		t.AssertNil(err)
		t.Assert(s.GetListenedPort() > 0, true)

		err = s.StartManaged()
		t.AssertNE(err, nil)

		s.StopForceful()
	})
}

func Test_Grpcx_GappServerAdapter_StartGracefulStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := grpcx.Server.NewConfig()
		c.Name = guid.S()
		s := grpcx.Server.New(c)
		controller.Register(s)

		ad := grpcx.NewGappServerAdapter(s)

		err := ad.Start()
		t.AssertNil(err)
		t.Assert(s.GetListenedPort() > 0, true)

		err = ad.Stop(true)
		t.AssertNil(err)
	})
}

func Test_Grpcx_GappServerAdapter_StopForceful(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := grpcx.Server.NewConfig()
		c.Name = guid.S()
		s := grpcx.Server.New(c)
		controller.Register(s)

		ad := grpcx.NewGappServerAdapter(s)

		err := ad.Start()
		t.AssertNil(err)
		t.Assert(s.GetListenedPort() > 0, true)

		err = ad.Stop(false)
		t.AssertNil(err)
	})
}
