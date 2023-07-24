# nacos

Package `nacos` implements GoFrame `gcfg.Adapter` using nacos service.

# Installation

```
go get -u github.com/gogf/gf/contrib/config/nacos/v2
```

# Usage

## Create a custom boot package

If you wish using configuration from nacos globally,
it is strongly recommended creating a custom boot package in very top import,
which sets the Adapter of default configuration instance before any other package boots.

```go
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
	// Create anacosClient that implements gcfg.Adapter.
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
```

## Import boot package in top of main

It is strongly recommended import your boot package in top of your `main.go`.

Note the top `import`: `_ "github.com/gogf/gf/example/config/nacos/boot"` .

```go
package main

import (
    _ "github.com/gogf/gf/example/config/nacos/boot"

    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    var ctx = gctx.GetInitCtx()

    // Available checks.
    g.Dump(g.Cfg().Available(ctx))

    // All key-value configurations.
    g.Dump(g.Cfg().Data(ctx))

    // Retrieve certain value by key.
    g.Dump(g.Cfg().MustGet(ctx, "server.address"))
}
```

## License

`GoFrame nacos` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.
