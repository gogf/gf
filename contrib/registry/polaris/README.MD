# GoFrame Polaris Registry

English | [简体中文](README_ZH.MD)

Use `PolarisMesh` as service registration, discovery management and heartbeat reporting.

## Installation

```
go get -u -v github.com/gogf/gf/contrib/registry/polaris/v2
```

suggested using `go.mod`:

```
require github.com/gogf/gf/contrib/registry/polaris/v2 latest
```

## Example

### Reference example

[server](example/registry/polaris/server/main.go)

```go
package main

import (
    "github.com/polarismesh/polaris-go/pkg/config"
    
    "github.com/gogf/gf/contrib/registry/polaris/v2"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/net/gsvc"
)

func main() {
    conf := config.NewDefaultConfiguration([]string{"192.168.100.222:8091"})
    
    // TTL egt 2*time.Second
    gsvc.SetRegistry(polaris.NewRegistryWithConfig(conf, polaris.WithTTL(100)))
    
    s := g.Server(`hello.svc`)
    s.BindHandler("/", func(r *ghttp.Request) {
        g.Log().Info(r.Context(), `request received`)
        r.Response.Write(`Hello world`)
    })
    s.Run()
}
```

[client](example/registry/polaris/client/main.go)

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/polarismesh/polaris-go/pkg/config"
    
    "github.com/gogf/gf/contrib/registry/polaris/v2"
    
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/gsvc"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    conf := config.NewDefaultConfiguration([]string{"192.168.100.222:8091"})
    
    gsvc.SetRegistry(polaris.NewRegistryWithConfig(conf, polaris.WithTTL(100)))
    
    for i := 0; i < 100; i++ {
        res, err := g.Client().Get(gctx.New(), `http://hello.svc/`)
        if err != nil {
            panic(err)
        }
        fmt.Println(res.ReadAllString())
        res.Close()
        time.Sleep(time.Second)
    }
}
```

## License

`GoFrame Polaris` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.

