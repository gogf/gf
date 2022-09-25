package boot

import (
	"github.com/gogf/gf/contrib/config/kubecm/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	configmapName       = "test-configmap"
	dataItemInConfigmap = "config.yaml"
)

func init() {
	var (
		err error
		ctx = gctx.GetInitCtx()
	)
	// Create kubecm Client that implements gcfg.Adapter.
	adapter, err := kubecm.New(gctx.GetInitCtx(), kubecm.Config{
		ConfigMap: configmapName,
		DataItem:  dataItemInConfigmap,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}
