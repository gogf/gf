package main

import (
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	gsvc.SetRegistry(etcd.New(`127.0.0.1:2379`))
	gsel.SetBuilder(gsel.NewBuilderRoundRobin())

	for i := 0; i < 10; i++ {
		ctx := gctx.New()
		res := g.Client().GetContent(ctx, `http://hello.svc/`)
		g.Log().Info(ctx, res)
	}
}
