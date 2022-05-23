package main

import (
	"context"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/contrib/registry/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
)

func main() {
	conf := config.NewDefaultConfiguration([]string{"192.168.100.222:8091"})
	conf.Consumer.LocalCache.SetPersistDir("/tmp/polaris/backup")
	err := api.SetLoggersDir("/tmp/polaris/log")
	if err != nil {
		g.Log().Fatal(context.Background(), err)
	}

	// TTL egt 2*time.Second
	gsvc.SetRegistry(polaris.NewWithConfig(conf, polaris.WithTTL(100)))

	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
