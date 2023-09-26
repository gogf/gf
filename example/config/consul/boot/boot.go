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
