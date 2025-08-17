# GoFrame Consul Registry

Use `consul` as service registration and discovery management.

## Installation
```bash
go get -u github.com/gogf/gf/contrib/registry/consul/v2
```
suggested using `go.mod`:
```bash
require github.com/gogf/gf/contrib/registry/consul/v2 latest
```

## Example

### HTTP Server
```go
package main

import (
    "context"
    
    "github.com/gogf/gf/contrib/registry/consul/v2"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
    "github.com/gogf/gf/v2/net/gsvc"
)

func main() {
    registry, err := consul.New(consul.WithAddress("127.0.0.1:8500"))
    if err != nil {
        g.Log().Fatal(context.Background(), err)
    }
    gsvc.SetRegistry(registry)

    s := g.Server("hello.svc")
    s.BindHandler("/", func(r *ghttp.Request) {
        g.Log().Info(r.Context(), "request received")
        r.Response.Write("Hello world")
    })
    s.Run()
}
```

### HTTP Client
```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/gogf/gf/contrib/registry/consul/v2"
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/gsel"
    "github.com/gogf/gf/v2/net/gsvc"
    "github.com/gogf/gf/v2/os/gctx"
)

func main() {
    registry, err := consul.New(consul.WithAddress("127.0.0.1:8500"))
    if err != nil {
        g.Log().Fatal(context.Background(), err)
    }
    gsvc.SetRegistry(registry)
    gsel.SetBuilder(gsel.NewBuilderRoundRobin())

    client := g.Client()
    for i := 0; i < 100; i++ {
        res, err := client.Get(gctx.New(), "http://hello.svc/")
        if err != nil {
            panic(err)
        }
        fmt.Println(res.ReadAllString())
        res.Close()
        time.Sleep(time.Second)
    }
}
```

## Configuration Options

The registry supports the following configuration options:

- `WithAddress(address string)`: Sets the Consul server address (default: "127.0.0.1:8500")
- `WithToken(token string)`: Sets the ACL token for Consul authentication

## Features

- Service registration with TTL health check
- Service discovery with health status filtering
- Service metadata support
- Watch support for service changes
- Consul ACL token support

## Requirements

- Go 1.18 or higher
- Consul 1.0 or higher

## License

`GoFrame Consul` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.
