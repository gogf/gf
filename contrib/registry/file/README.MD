# GoFrame File Registry


Use `file` as service registration and discovery management.


## Installation
```
go get -u -v github.com/gogf/gf/contrib/registry/file/v2
```
suggested using `go.mod`:
```
require github.com/gogf/gf/contrib/registry/file/v2 latest
```


## Example

### Reference example

[server](example/registry/file/server/main.go)
```go
package main

import (
	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	gsvc.SetRegistry(file.New(gfile.Temp("gsvc")))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}

```

[client](example/registry/file/client/main.go)
```go
package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

func main() {
	gsvc.SetRegistry(file.New(gfile.Temp("gsvc")))

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

