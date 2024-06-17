// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package boot

import (
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/gogf/gf/contrib/config/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
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
	// Create nacos Client that implements gcfg.Adapter.
	adapter, err := nacos.New(ctx, nacos.Config{
		ServerConfigs: []constant.ServerConfig{serverConfig},
		ClientConfig:  clientConfig,
		ConfigParam:   configParam,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}
