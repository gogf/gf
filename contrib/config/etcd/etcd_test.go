// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd_test

import (
	"github.com/gogf/gf/contrib/config/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
)

var (
	ctx     = gctx.GetInitCtx()
	Endpoints   = []string{"localhost:2379"}
	ConfigKey = "/configs/config.yaml"
	Watch = false
)

func TestEtcd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, err := etcd.New(ctx, etcd.Config{
			Endpoints:   Endpoints,
			ConfigKey: ConfigKey,
			Watch: Watch,
		})
		t.AssertNil(err)

		config := g.Cfg(guid.S())
		config.SetAdapter(adapter)
		t.Assert(config.Available(ctx), true)

		v, err := config.Get(ctx, `server.address`)
		t.AssertNil(err)
		t.Assert(v.String(), ":8000")

		m, err := config.Data(ctx)
		t.AssertNil(err)
		t.AssertGT(len(m), 0)
	})
}
