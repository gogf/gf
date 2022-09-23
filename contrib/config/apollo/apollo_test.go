package apollo

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/spf13/viper"
	"testing"
)

func TestApollo(t *testing.T) {
	ctx := gctx.New()

	NewApollo(appId, cluster, ip).Run()

	g.Dump(g.Cfg().Data(ctx))
	g.Dump(viper.AllSettings())
}
