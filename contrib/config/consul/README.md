# consul

Package `consul` implements GoFrame `gcfg.Adapter` using consul service.

# Installation

```
go get -u github.com/gogf/gf/contrib/config/consul/v2
```

# Usage

## Create a custom boot package

If you wish using configuration from consul globally,
it is strongly recommended creating a custom boot package in very top import,
which sets the Adapter of default configuration instance before any other package boots.

```go
package boot

import (
	consul "github.com/gogf/gf/contrib/config/consul/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

func init() {
	var (
		ctx          = gctx.GetInitCtx()
		consulConfig = api.Config{
			Address:    "127.0.0.1:8500",
			Scheme:     "http",
			Datacenter: "dc1",
			Transport:  cleanhttp.DefaultPooledTransport(),
			Token:      "3f8aeba2-f1f7-42d0-b912-fcb041d4546d",
		}
		configPath = "server/message"
	)

	adapter, err := consul.New(ctx, consul.Config{
		ConsulConfig: consulConfig,
		Path:   configPath,
		Watch:  true,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `New consul adapter error: %+v`, err)
	}

	g.Cfg().SetAdapter(adapter)
}
```

## Import boot package in top of main

It is strongly recommended import your boot package in top of your `main.go`.

Note the top `import`: `_ "github.com/gogf/gf/example/config/consul/boot"` .

```go
package main

import (
	_ "github.com/gogf/gf/example/config/consul/boot"

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
	g.Dump(g.Cfg().MustGet(ctx, "redis.addr"))
}
```

## License

`GoFrame consul` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.