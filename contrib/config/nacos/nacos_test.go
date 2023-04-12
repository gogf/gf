// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/gogf/gf/contrib/config/nacos/v2"
)

var (
	ctx          = gctx.GetInitCtx()
	serverConfig = constant.ServerConfig{
		IpAddr: "localhost",
		Port:   8848,
	}
	clientConfig = constant.ClientConfig{
		CacheDir: "/tmp/nacos",
		LogDir:   "/tmp/nacos",
	}
	configParam = vo.ConfigParam{
		DataId: "config.toml",
		Group:  "test",
	}
)

func TestNacos(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, err := nacos.New(ctx, nacos.Config{
			ServerConfigs: []constant.ServerConfig{serverConfig},
			ClientConfig:  clientConfig,
			ConfigParam:   configParam,
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
