package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	_ "github.com/gogf/gf/example/config/kubecm/boot_in_pod"
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
