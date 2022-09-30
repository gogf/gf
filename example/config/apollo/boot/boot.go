package boot

import (
	"github.com/gogf/gf/contrib/config/apollo/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	var (
		ctx     = gctx.GetInitCtx()
		appId   = "SampleApp"
		cluster = "default"
		ip      = "http://localhost:8080"
	)
	// Create apollo Client that implements gcfg.Adapter.
	adapter, err := apollo.New(ctx, apollo.Config{
		AppID:   appId,
		IP:      ip,
		Cluster: cluster,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}
