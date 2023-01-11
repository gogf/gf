# GoFrame Etcd Registry

Use `zookeeper` as service registration and discovery management.

## Installation
```
go get -u -v github.com/gogf/gf/contrib/registry/zookeeper/v2
```
suggested using `go.mod`:
```
require github.com/gogf/gf/contrib/registry/zookeeper/v2 latest
```

```go
package main

import (
	"github.com/gogf/gf/contrib/registry/zookeeper/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
)

func main() {
	gsvc.SetRegistry(zookeeper.New(
		[]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogf"),
    ))
	
	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
```

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/zookeeper/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	gsvc.SetRegistry(zookeeper.New(
		[]string{"127.0.0.1:2181"}, zookeeper.WithRootPath("/gogf"),
	))
	gsel.SetBuilder(gsel.NewBuilderRoundRobin())

	client := g.Client()
	for i := 0; i < 100; i++ {
		res, err := client.Get(gctx.New(), `http://hello.svc/`)
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

`GoFrame zookeeper` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.

