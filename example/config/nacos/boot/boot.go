package boot

import (
	"github.com/gogf/gf/contrib/config/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
