package main

import (
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
	for i := 0; i < 10; i++ {
		ctx := gctx.New()
		res, err := client.Get(ctx, `http://hello.svc/`)
		if err != nil {
			panic(err)
		}
		g.Log().Debug(ctx, res.ReadAllString())
		res.Close()
		time.Sleep(time.Second)
	}
}
