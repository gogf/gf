# GoFrame Nacos Registry


Use `nacos` as service registration and discovery management.


## Installation
```
go get -u -v github.com/wangyougui/gf/contrib/registry/nacos/v2
```
suggested using `go.mod`:
```
require github.com/wangyougui/gf/contrib/registry/nacos/v2 latest
```


## Example

### Reference example

[server](example/registry/nacos/server/main.go)
```go
package main

import (
	"github.com/wangyougui/gf/contrib/registry/nacos/v2"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/net/ghttp"
	"github.com/wangyougui/gf/v2/net/gsvc"
)

func main() {
	gsvc.SetRegistry(nacos.New(`127.0.0.1:8848`).
        SetClusterName("DEFAULT").
        SetGroupName("DEFAULT_GROUP"))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
```

[client](example/registry/nacos/client/main.go)
```go
package main

import (
	"fmt"
	"time"

	"github.com/wangyougui/gf/contrib/registry/nacos/v2"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/net/gsel"
	"github.com/wangyougui/gf/v2/net/gsvc"
	"github.com/wangyougui/gf/v2/os/gctx"
)

func main() {
	gsvc.SetRegistry(nacos.New(`127.0.0.1:8848`))
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

`GoFrame Nacos` is licensed under the [MIT License](../../../LICENSE), 100% free and open-source, forever.

